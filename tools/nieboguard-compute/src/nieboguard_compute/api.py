"""HTTPS API client with strict URL, size and checksum checks."""
from __future__ import annotations

import json
import os
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import Any
from urllib.parse import urlparse

from . import __version__
from .crypto import sha256_hex

MAX_JSON = 2 * 1024 * 1024
LOCAL_HOSTS = {"localhost", "127.0.0.1", "::1"}


class ApiError(RuntimeError):
    def __init__(self, status: int, code: str, message: str, data: Any = None):
        super().__init__(f"{status} {code}: {message}")
        self.status = status
        self.code = code
        self.message = message
        self.data = data


class _NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):  # type: ignore[override]
        raise urllib.error.HTTPError(req.full_url, code, "redirect refused", headers, fp)


_OPENER = urllib.request.build_opener(_NoRedirect)


def _require_https(url: str, allow_http_local: bool = True) -> None:
    parsed = urlparse(url)
    if parsed.scheme == "https":
        return
    if allow_http_local and parsed.scheme == "http" and parsed.hostname in LOCAL_HOSTS:
        return
    raise ApiError(0, "insecure_url", "Dozwolone są wyłącznie adresy HTTPS.")


def same_origin(url: str, api_base: str) -> bool:
    a = urlparse(url)
    b = urlparse(api_base)
    return (a.scheme, a.hostname, a.port) == (b.scheme, b.hostname, b.port)


class Client:
    def __init__(self, api_base: str, access_token: str | None = None, timeout: float = 30.0):
        self.api_base = api_base.rstrip("/")
        _require_https(self.api_base)
        self.access_token = access_token
        self.timeout = timeout

    def _request(self, method: str, path: str, body: Any = None, token: str | None = None,
                 expected: tuple[int, ...] = (200,)) -> Any:
        url = self.api_base + path
        _require_https(url)
        data = None
        headers = {
            "Accept": "application/json",
            "User-Agent": f"NieboGuard-Compute/{__version__}",
        }
        if body is not None:
            data = json.dumps(body, ensure_ascii=False).encode("utf-8")
            headers["Content-Type"] = "application/json"
        bearer = token if token is not None else self.access_token
        if bearer:
            headers["Authorization"] = f"Bearer {bearer}"
        req = urllib.request.Request(url, data=data, headers=headers, method=method)
        try:
            with _OPENER.open(req, timeout=self.timeout) as resp:
                raw = resp.read(MAX_JSON + 1)
                status = resp.status
                if len(raw) > MAX_JSON:
                    raise ApiError(status, "too_large", "Odpowiedź API jest zbyt duża.")
                payload = json.loads(raw.decode("utf-8")) if raw else None
        except urllib.error.HTTPError as exc:
            raw = exc.read(MAX_JSON) if exc.fp else b""
            try:
                payload = json.loads(raw.decode("utf-8")) if raw else None
            except Exception:
                payload = None
            error = (payload or {}).get("error") if isinstance(payload, dict) else None
            code = (error or {}).get("code") if isinstance(error, dict) else "http_error"
            message = (error or {}).get("message") if isinstance(error, dict) else f"HTTP {exc.code}"
            raise ApiError(exc.code, str(code), str(message), payload) from exc
        except urllib.error.URLError as exc:
            raise ApiError(0, "network_error", f"Brak połączenia z API: {exc.reason}") from exc
        if status not in expected:
            raise ApiError(status, "unexpected_status", f"Nieoczekiwany status HTTP {status}.", payload)
        return payload

    # ------------------------------------------------------------------
    def status(self) -> dict[str, Any]:
        return self._request("GET", "/api/v2/compute/status", expected=(200,))

    def leaderboard(self, limit: int = 50) -> dict[str, Any]:
        return self._request("GET", f"/api/v2/compute/leaderboard?limit={int(limit)}", expected=(200,))

    def enroll(self, invite_code: str, hardware: dict[str, Any], consent_version: str,
               preferences: dict[str, Any], display_name: str | None,
               leaderboard_opt_in: bool, capabilities: dict[str, Any]) -> dict[str, Any]:
        return self._request("POST", "/api/v2/compute/enroll", {
            "invite_code": invite_code,
            "accept_consent": True,
            "consent_version": consent_version,
            "hardware": {k: hardware.get(k) for k in (
                "os_family", "arch", "cpu_cores", "ram_gb", "gpu_vendor", "gpu_model", "gpu_memory_gb")},
            "preferences": preferences,
            "display_name": display_name,
            "leaderboard_opt_in": bool(leaderboard_opt_in),
            "capabilities": capabilities,
            "client": {"version": __version__, "protocol": 1},
        }, expected=(201,))

    def refresh(self, device_id: str, refresh_token: str) -> dict[str, Any]:
        data = self._request("POST", "/api/v2/compute/token", {
            "device_id": device_id, "refresh_token": refresh_token,
        }, expected=(200,))
        self.access_token = data.get("access_token")
        return data

    def me(self) -> dict[str, Any]:
        return self._request("GET", "/api/v2/compute/me", expected=(200,))

    def update_preferences(self, preferences: dict[str, Any]) -> dict[str, Any]:
        return self._request("PUT", "/api/v2/compute/me/preferences", {"preferences": preferences}, expected=(200,))

    def revoke(self, reason: str) -> dict[str, Any]:
        return self._request("POST", "/api/v2/compute/consent/revoke", {"reason": reason}, expected=(200,))

    def next_task(self) -> dict[str, Any]:
        return self._request("GET", "/api/v2/compute/tasks/next", expected=(200,))

    def heartbeat(self, task_id: str, progress: float) -> dict[str, Any]:
        return self._request("POST", f"/api/v2/compute/tasks/{urllib.parse.quote(task_id)}/heartbeat",
                             {"progress": progress}, expected=(200,))

    def fail(self, task_id: str, error: str) -> dict[str, Any]:
        return self._request("POST", f"/api/v2/compute/tasks/{urllib.parse.quote(task_id)}/fail",
                             {"error": error[:300]}, expected=(200,))

    def complete(self, task_id: str, payload: dict[str, Any]) -> dict[str, Any]:
        return self._request("POST", f"/api/v2/compute/tasks/{urllib.parse.quote(task_id)}/complete",
                             payload, expected=(200, 202, 422))

    # ------------------------------------------------------------------
    def download(self, url: str, dest: Path, max_bytes: int, sha256: str | None = None) -> Path:
        _require_https(url)
        if not same_origin(url, self.api_base):
            raise ApiError(0, "untrusted_origin", "Adres danych nie należy do skonfigurowanego API.")
        if max_bytes <= 0 or max_bytes > 512 * 1024 * 1024:
            raise ApiError(0, "bad_limit", "Nieprawidłowy limit pobierania.")
        dest.parent.mkdir(parents=True, exist_ok=True)
        tmp = dest.with_name(dest.name + ".part")
        digest = __import__("hashlib").sha256()
        total = 0
        req = urllib.request.Request(url, headers={"User-Agent": f"NieboGuard-Compute/{__version__}"})
        try:
            with _OPENER.open(req, timeout=self.timeout) as resp:
                if resp.status != 200:
                    raise ApiError(resp.status, "download_failed", f"HTTP {resp.status} przy pobieraniu.")
                with open(tmp, "wb") as fh:
                    while True:
                        chunk = resp.read(1024 * 256)
                        if not chunk:
                            break
                        total += len(chunk)
                        if total > max_bytes:
                            raise ApiError(0, "too_large", "Pobierany plik przekracza limit.")
                        digest.update(chunk)
                        fh.write(chunk)
        except urllib.error.HTTPError as exc:
            raise ApiError(exc.code, "download_failed", f"HTTP {exc.code} przy pobieraniu.") from exc
        except urllib.error.URLError as exc:
            raise ApiError(0, "network_error", f"Błąd pobierania: {exc.reason}") from exc
        actual = digest.hexdigest()
        if sha256 and actual != sha256.lower():
            tmp.unlink(missing_ok=True)
            raise ApiError(0, "checksum_mismatch", "Suma SHA-256 pobranego pliku jest niezgodna.")
        os.replace(tmp, dest)
        return dest

    def upload(self, url: str, data: bytes, content_type: str, sha256: str) -> dict[str, Any]:
        _require_https(url)
        if not same_origin(url, self.api_base):
            raise ApiError(0, "untrusted_origin", "Adres uploadu nie należy do skonfigurowanego API.")
        if len(data) > 16 * 1024 * 1024:
            raise ApiError(0, "too_large", "Artefakt przekracza limit klienta.")
        req = urllib.request.Request(url, data=data, method="PUT", headers={
            "Content-Type": content_type,
            "Content-Length": str(len(data)),
            "X-Content-SHA256": sha256,
            "User-Agent": f"NieboGuard-Compute/{__version__}",
        })
        try:
            with _OPENER.open(req, timeout=max(self.timeout, 120)) as resp:
                raw = resp.read(MAX_JSON)
                return json.loads(raw.decode("utf-8")) if raw else {}
        except urllib.error.HTTPError as exc:
            raw = exc.read(MAX_JSON) if exc.fp else b""
            try:
                payload = json.loads(raw.decode("utf-8")) if raw else None
            except Exception:
                payload = None
            error = (payload or {}).get("error") if isinstance(payload, dict) else None
            code = (error or {}).get("code") if isinstance(error, dict) else "upload_failed"
            message = (error or {}).get("message") if isinstance(error, dict) else f"HTTP {exc.code}"
            raise ApiError(exc.code, str(code), str(message), payload) from exc
        except urllib.error.URLError as exc:
            raise ApiError(0, "network_error", f"Błąd uploadu: {exc.reason}") from exc


def retry(call, attempts: int = 4, base_delay: float = 1.0):
    last = None
    for i in range(attempts):
        try:
            return call()
        except ApiError as exc:
            last = exc
            if exc.status not in (429, 500, 502, 503, 504, 0):
                raise
            time.sleep(min(30.0, base_delay * (2 ** i)))
    if last:
        raise last
