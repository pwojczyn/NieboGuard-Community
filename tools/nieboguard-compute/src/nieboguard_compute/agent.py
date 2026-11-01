"""Opcjonalny lokalny agent dla strony /compute/.

Agent działa wyłącznie na 127.0.0.1, uruchamia się ręcznie i nie startuje
treningu bez żądania POST. Jest to jedyny sposób, w jaki przeglądarka może
sterować lokalnym klientem. Nie wystawiaj go na LAN.
"""
from __future__ import annotations

import json
import os
import subprocess
import sys
import threading
import time
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from typing import Any

from .config import data_root, load

ALLOWED_ORIGINS = {
    "https://nieboguard.pl", "https://www.nieboguard.pl",
    "http://localhost:8788", "http://127.0.0.1:8788",
    "http://localhost:3000", "http://127.0.0.1:3000",
}
MAX_BODY = 4096


def _pid_alive(pid: int) -> bool:
    if pid <= 0:
        return False
    try:
        if os.name == "nt":
            out = subprocess.run(["tasklist", "/FI", f"PID eq {pid}"], capture_output=True, text=True, timeout=5)
            return str(pid) in out.stdout
        os.kill(pid, 0)
        return True
    except Exception:
        return False


class AgentState:
    def __init__(self, cfg: dict[str, Any], host: str = "127.0.0.1", port: int = 8765):
        self.cfg = cfg
        self.root = data_root(cfg)
        self.pid_file = self.root / "agent-runner.pid"
        self.stop_file = self.root / "STOP"
        self.log_file = self.root / "agent-runner.log"
        self.host = host
        self.port = port
        self.lock = threading.Lock()

    def runner_pid(self) -> int:
        try:
            pid = int(self.pid_file.read_text(encoding="utf-8").strip())
        except Exception:
            return 0
        if not _pid_alive(pid):
            self.pid_file.unlink(missing_ok=True)
            return 0
        return pid

    def status(self) -> dict[str, Any]:
        pid = self.runner_pid()
        return {
            "ok": True,
            "agent": "nieboguard-compute",
            "runner_running": bool(pid),
            "runner_pid": pid or None,
            "stop_requested": self.stop_file.exists(),
            "enrolled": bool(self.cfg.get("device_id") and self.cfg.get("refresh_token")),
            "schedule": self.cfg.get("schedule"),
            "api_base": self.cfg.get("api_base"),
            "data_dir": str(self.root),
        }

    def start(self) -> tuple[int, dict[str, Any]]:
        with self.lock:
            if not self.cfg.get("device_id") or not self.cfg.get("refresh_token"):
                return 409, {"error": {"code": "not_enrolled", "message": "Najpierw zarejestruj urządzenie."}}
            if self.runner_pid():
                return 409, {"error": {"code": "already_running", "message": "Klient już działa."}}
            if not self.cfg.get("consent_version"):
                return 409, {"error": {"code": "consent_required", "message": "Brak zapisanej zgody."}}
            self.stop_file.unlink(missing_ok=True)
            env = dict(os.environ)
            src = str(Path(__file__).resolve().parents[1])
            env["PYTHONPATH"] = src + (os.pathsep + env["PYTHONPATH"] if env.get("PYTHONPATH") else "")
            log = open(self.log_file, "ab", buffering=0)
            try:
                kwargs: dict[str, Any] = {"stdout": log, "stderr": log, "stdin": subprocess.DEVNULL, "env": env}
                if os.name == "nt":
                    kwargs["creationflags"] = getattr(subprocess, "CREATE_NEW_PROCESS_GROUP", 0)
                else:
                    kwargs["start_new_session"] = True
                proc = subprocess.Popen([sys.executable, "-m", "nieboguard_compute", "--data-dir", str(self.root), "start"], **kwargs)
            finally:
                log.close()
            self.pid_file.write_text(str(proc.pid), encoding="utf-8")
            return 202, {"ok": True, "runner_pid": proc.pid, "note": "Klient uruchomiony na Twoją komendę. Postęp w logu lokalnym."}

    def stop(self) -> tuple[int, dict[str, Any]]:
        self.stop_file.write_text("stop\n", encoding="utf-8")
        return 200, {"ok": True, "stop_requested": True, "note": "Zatrzymanie po checkpointcie; proces nie jest zabijany siłą."}

    def resume(self) -> tuple[int, dict[str, Any]]:
        self.stop_file.unlink(missing_ok=True)
        return 200, {"ok": True, "stop_requested": False}


def _handler(state: AgentState):
    class Handler(BaseHTTPRequestHandler):
        server_version = "NieboGuardComputeAgent/0.1"
        def log_message(self, fmt, *args):  # noqa: A003
            return
        def _origin(self) -> str | None:
            origin = self.headers.get("Origin")
            return origin if origin in ALLOWED_ORIGINS else None
        def _cors(self):
            origin = self._origin()
            if origin:
                self.send_header("Access-Control-Allow-Origin", origin)
                self.send_header("Vary", "Origin")
                self.send_header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
                self.send_header("Access-Control-Allow-Headers", "Content-Type, X-NieboGuard-Agent")
        def _json(self, status: int, payload: dict[str, Any]):
            body = json.dumps(payload, ensure_ascii=False).encode()
            self.send_response(status)
            self.send_header("Content-Type", "application/json; charset=utf-8")
            self.send_header("Content-Length", str(len(body)))
            self.send_header("Cache-Control", "no-store")
            self.send_header("X-Content-Type-Options", "nosniff")
            self._cors()
            self.end_headers()
            self.wfile.write(body)
        def do_OPTIONS(self):
            if not self._origin():
                self._json(403, {"error": {"code": "origin_denied", "message": "Origin niedozwolony."}})
                return
            self.send_response(204)
            self._cors()
            self.end_headers()
        def do_GET(self):
            if self.headers.get("Origin") and not self._origin():
                self._json(403, {"error": {"code": "origin_denied", "message": "Origin niedozwolony."}})
                return
            if self.path.split("?")[0] not in ("/status", "/health"):
                self._json(404, {"error": {"code": "not_found", "message": "Nie znaleziono."}})
                return
            self._json(200, state.status())
        def do_POST(self):
            origin = self._origin()
            # A browser request must carry an allowlisted Origin. A local
            # command-line request may omit it, which is useful for testing.
            if self.headers.get("Origin") and not origin:
                self._json(403, {"error": {"code": "origin_denied", "message": "Origin niedozwolony."}})
                return
            path = self.path.split("?")[0]
            if path not in ("/start", "/stop", "/resume"):
                self._json(404, {"error": {"code": "not_found", "message": "Nie znaleziono."}})
                return
            length = int(self.headers.get("Content-Length") or 0)
            if length > MAX_BODY:
                self._json(413, {"error": {"code": "too_large", "message": "Ciało żądania jest zbyt duże."}})
                return
            if length:
                self.rfile.read(length)
            if path == "/start":
                status, payload = state.start()
            elif path == "/stop":
                status, payload = state.stop()
            else:
                status, payload = state.resume()
            self._json(status, payload)
    return Handler


def serve(cfg: dict[str, Any], host: str = "127.0.0.1", port: int = 8765) -> None:
    state = AgentState(cfg, host, port)
    server = ThreadingHTTPServer((host, port), _handler(state))
    print(f"NieboGuard Compute agent: http://{host}:{port}")
    print("Agent działa tylko lokalnie. Otwórz https://nieboguard.pl/compute/ i użyj przycisków.")
    print("Zatrzymanie agenta: Ctrl+C. To nie zatrzymuje treningu; użyj przycisku Stop.")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        pass
    finally:
        server.server_close()


def main(cfg: dict[str, Any], host: str = "127.0.0.1", port: int = 8765) -> int:
    serve(cfg, host, port)
    return 0
