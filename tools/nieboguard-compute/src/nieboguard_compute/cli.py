"""Command line interface for NieboGuard Compute."""
from __future__ import annotations

import argparse
import json
import os
import platform
import sys
import time
from pathlib import Path
from typing import Any

from . import __version__
from .api import ApiError, Client
from .config import (CONSENT_VERSION, ConfigError, config_path, data_root, default_data_dir,
                     ensure_data_dir, load, redacted, save, validate)
from .hardware import detect, idle_seconds, summary
from .runner import Runner, RunnerError
from .torch_backend import describe as torch_describe
from .agent import serve as serve_agent

CONSENT_TEXT = """\
NieboGuard Compute wykorzystuje Twój komputer do lokalnego trenowania małych
modeli AI. Udział jest dobrowolny i wymaga tej zgody. Zasady:

- program nie uruchamia się ukrycie; startujesz go samodzielnie,
- dane treningowe są publiczne i mają potwierdzone prawo redystrybucji,
- Twój komputer wysyła wynikowe wagi, metryki i minimalne dane sprzętowe
  (system, architektura, liczba rdzeni, kubełek RAM, model GPU),
- nie wysyłamy nazwy hosta, adresów, plików użytkownika ani prywatnych zdjęć,
- możesz ograniczyć CPU/GPU, wybrać harmonogram i w każdej chwili cofnąć dostęp,
- trening zużywa prąd i obciąża sprzęt; laptop na baterii może pracować krócej,
- wynik trafia do kwarantanny i niezależnej ewaluacji; nie jest publikowany
  automatycznie,
- program nie pobiera ani nie wykonuje kodu z backendu; manifest jest podpisany,
- w razie wątpliwości zatrzymaj program i usuń konfigurację.

Nie jest to system ostrzegania ludności. Nie wysyłaj danych osobowych.
"""


def _print(obj: Any) -> None:
    print(json.dumps(obj, ensure_ascii=False, indent=2, sort_keys=True))


def _load(args) -> dict[str, Any]:
    return load(getattr(args, "config", None), getattr(args, "data_dir", None))


def _save(cfg: dict[str, Any], args) -> None:
    save(cfg, getattr(args, "config", None))


def _client(cfg: dict[str, Any]) -> Client:
    return Client(str(cfg["api_base"]))


def cmd_doctor(args) -> int:
    hw = detect()
    cfg = _load(args)
    result: dict[str, Any] = {
        "client_version": __version__,
        "python": sys.version.split()[0],
        "platform": platform.platform(),
        "hardware": hw,
        "hardware_summary": summary(hw),
        "data_dir": str(data_root(cfg)),
        "config_path": str(config_path(cfg.get("data_dir"))),
        "configured": bool(cfg.get("device_id")),
        "manifest_key_pinned": bool(cfg.get("manifest_public_key")),
        "consent_version": cfg.get("consent_version"),
        "supported": hw.get("supported"),
        "torch": torch_describe(),
    }
    try:
        status = _client(cfg).status()
        result["backend"] = {
            "ok": True, "enabled": status.get("enabled"), "configured": status.get("configured"),
            "consent_version": (status.get("consent") or {}).get("version"),
            "participants": status.get("participants"), "active_24h": status.get("active_24h"),
        }
        if not cfg.get("manifest_public_key") and status.get("manifest_public_key"):
            result["backend"]["manifest_public_key_available"] = True
    except ApiError as exc:
        result["backend"] = {"ok": False, "error": exc.message}
    idle = idle_seconds()
    result["idle_supported"] = idle is not None
    if args.json:
        _print(result)
    else:
        print("NieboGuard Compute — diagnostyka")
        print(result["hardware_summary"])
        print(f"Python: {result['python']} · dane: {result['data_dir']}")
        print(f"Backend: {'OK' if result['backend'].get('ok') else 'BŁĄD'} · "
              f"{result['backend'].get('error', 'skonfigurowany=' + str(result['backend'].get('configured')))}")
        print(f"Rejestracja: {'tak' if result['configured'] else 'nie'} · "
              f"klucz manifestu: {'przypięty' if result['manifest_key_pinned'] else 'brak'}")
        print(f"Tryb bezczynności: {'dostępny' if result['idle_supported'] else 'niedostępny (użyj manual/always)'}")
        t = result.get("torch") or {}
        print(f"PyTorch: {'tak' if t.get('available') else 'nie (opcjonalny)'} · CUDA: {'tak' if t.get('cuda') else 'nie'}")
        if not hw.get("supported"):
            print("UWAGA: ten system/architektura nie jest jeszcze obsługiwana przez backend.")
    return 0 if hw.get("supported") and result.get("backend", {}).get("ok") else 1


def cmd_enroll(args) -> int:
    if not args.consent:
        print("Wymagana jawna zgoda. Przeczytaj tekst i dodaj flagę --consent.\n")
        print(CONSENT_TEXT)
        return 2
    print(CONSENT_TEXT)
    cfg = _load(args)
    client = _client(cfg)
    try:
        status = client.status()
    except ApiError as exc:
        print(f"Backend niedostępny: {exc.message}")
        return 1
    if not status.get("enabled"):
        print("Backend NieboGuard Compute jest wyłączony.")
        return 1
    consent_version = str((status.get("consent") or {}).get("version") or CONSENT_VERSION)
    hw = detect()
    if not hw.get("supported"):
        print("Ten system/architektura nie jest obsługiwana (wymagane Windows/Linux x86_64).")
        return 1
    modalities = [args.modality] if args.modality in ("audio", "vision") else ["audio", "vision"]
    preferences = {
        "modalities": modalities,
        "max_cpu_percent": args.max_cpu,
        "max_gpu_percent": args.max_gpu,
        "schedule": args.schedule,
        "max_data_mb": args.max_data_mb,
        "pack_tiers": [x.strip().upper() for x in args.pack_tiers.split(",") if x.strip()],
    }
    try:
        result = client.enroll(
            invite_code=args.invite, hardware=hw, consent_version=consent_version,
            preferences=preferences, display_name=args.name,
            leaderboard_opt_in=args.leaderboard, capabilities={"backends": ["python"]},
        )
    except ApiError as exc:
        print(f"Rejestracja nie powiodła się: {exc.message}")
        return 1
    cfg.update({
        "device_id": result["device_id"], "refresh_token": result["refresh_token"],
        "consent_version": consent_version, "consent_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "display_name": args.name, "modalities": modalities,
        "max_cpu_percent": args.max_cpu, "max_gpu_percent": args.max_gpu,
        "schedule": args.schedule, "idle_minutes": args.idle_minutes,
        "max_data_mb": args.max_data_mb,
        "pack_tiers": [x.strip().upper() for x in args.pack_tiers.split(",") if x.strip()],
        "leaderboard_opt_in": bool(args.leaderboard),
        "manifest_public_key": result.get("manifest_public_key") or cfg.get("manifest_public_key") or "",
        "manifest_key_id": result.get("manifest_key_id") or cfg.get("manifest_key_id") or "ng-compute-1",
    })
    path = save(cfg)
    print(f"Zarejestrowano urządzenie {result['device_id']}.")
    print(f"Konfiguracja zapisana w {path} (tylko właściciel ma dostęp).")
    print("Refresh token nie jest wyświetlany. Aby uruchomić: nieboguard-compute start")
    print("Program nie uruchomi się w tle bez Twojej komendy.")
    return 0


def cmd_import_config(args) -> int:
    source = Path(args.file)
    try:
        imported = json.loads(source.read_text(encoding="utf-8"))
    except Exception as exc:
        print(f"Nie można odczytać pliku: {exc}")
        return 1
    if not isinstance(imported, dict) or not imported.get("device_id") or not imported.get("refresh_token"):
        print("Plik musi zawierać device_id i refresh_token.")
        return 1
    cfg = _load(args)
    for key in ("device_id", "refresh_token", "api_base", "manifest_public_key", "manifest_key_id", "consent_version", "display_name", "max_data_mb", "pack_tiers"):
        if imported.get(key):
            cfg[key] = imported[key]
    cfg["consent_at"] = cfg.get("consent_at") or time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime())
    try:
        path = save(cfg)
    except ConfigError as exc:
        print(f"Nieprawidłowa konfiguracja: {exc}")
        return 1
    print(f"Zaimportowano konfigurację urządzenia do {path}. Token nie jest wyświetlany.")
    return 0


def cmd_agent(args) -> int:
    cfg = _load(args)
    if not cfg.get("device_id") or not cfg.get("refresh_token"):
        print("Brak rejestracji. Najpierw enroll/import-config.")
        return 2
    if args.host not in ("127.0.0.1", "::1", "localhost"):
        print("Agent może nasłuchiwać wyłącznie na loopbacku.")
        return 2
    try:
        serve_agent(cfg, args.host, args.port)
    except OSError as exc:
        print(f"Nie można uruchomić agenta: {exc}")
        return 1
    return 0


def cmd_start(args) -> int:
    cfg = _load(args)
    if not cfg.get("device_id") or not cfg.get("refresh_token"):
        print("Brak rejestracji. Użyj: nieboguard-compute enroll --invite <kod> --consent")
        return 2
    if not cfg.get("consent_version"):
        print("Brak zapisanej zgody. Zarejestruj urządzenie ponownie.")
        return 2
    if not cfg.get("manifest_public_key") and not args.dev_allow_unsigned:
        print("Brak przypiętego klucza publicznego manifestu. Zaktualizuj konfigurację wydania "
              "albo użyj --dev-allow-unsigned wyłącznie w teście lokalnym.")
        return 2
    print("Uruchamiam za Twoją zgodą. Zatrzymanie: Ctrl+C lub nieboguard-compute stop. "
          "Program nie działa w tle.")
    stop = data_root(cfg) / "STOP"
    stop.unlink(missing_ok=True)
    runner = Runner(cfg, allow_unsigned=args.dev_allow_unsigned, verbose=True)
    try:
        return runner.run(once=args.once, max_tasks=args.max_tasks)
    except KeyboardInterrupt:
        stop.write_text("stop\n", encoding="utf-8")
        print("\nPrzerwano; zapisano żądanie zatrzymania.")
        return 130
    except (RunnerError, ConfigError) as exc:
        print(f"Nie można uruchomić: {exc}")
        return 1


def cmd_stop(args) -> int:
    cfg = _load(args)
    root = ensure_data_dir(cfg.get("data_dir"))
    (root / "STOP").write_text("stop\n", encoding="utf-8")
    print("Zapisano żądanie zatrzymania. Bieżące zadanie zakończy się po checkpointcie.")
    return 0


def cmd_resume(args) -> int:
    cfg = _load(args)
    (data_root(cfg) / "STOP").unlink(missing_ok=True)
    print("Usunięto żądanie zatrzymania. Uruchom: nieboguard-compute start")
    return 0


def cmd_status(args) -> int:
    cfg = _load(args)
    if not cfg.get("device_id") or not cfg.get("refresh_token"):
        print("Brak rejestracji.")
        return 2
    client = _client(cfg)
    try:
        client.refresh(str(cfg["device_id"]), str(cfg["refresh_token"]))
        data = client.me()
    except ApiError as exc:
        print(f"Błąd: {exc.message}")
        return 1
    if args.json:
        _print(data)
        return 0
    contribution = data.get("contribution") or {}
    device = data.get("device") or {}
    print(f"Urządzenie: {device.get('display_name') or device.get('id')} · {device.get('status')}")
    print(f"Zadania: {contribution.get('total_tasks', 0)} · zaakceptowane aktualizacje: {contribution.get('accepted_updates', 0)} · "
          f"czas: {int(contribution.get('total_seconds', 0)) // 60} min")
    print(f"Odznaki: {', '.join(contribution.get('badges') or []) or 'brak'}")
    for assignment in (data.get("assignments") or [])[:10]:
        print(f"  {assignment.get('assigned_at', '')} · {assignment.get('title')} · {assignment.get('status')} · {assignment.get('quality_status') or ''}")
    return 0


def cmd_revoke(args) -> int:
    cfg = _load(args)
    if not cfg.get("device_id") or not cfg.get("refresh_token"):
        print("Brak rejestracji.")
        return 2
    client = _client(cfg)
    try:
        client.refresh(str(cfg["device_id"]), str(cfg["refresh_token"]))
        data = client.revoke(args.reason)
    except ApiError as exc:
        print(f"Błąd: {exc.message}")
        return 1
    cfg["device_id"] = None
    cfg["refresh_token"] = None
    save(cfg)
    print("Dostęp cofnięty, tokeny unieważnione lokalnie i na backendzie." if data.get("revoked") else "Zgłoszono cofnięcie.")
    return 0


def cmd_config(args) -> int:
    cfg = _load(args)
    if args.action == "show":
        _print(redacted(cfg))
        return 0
    if args.action == "set":
        if not args.key or args.value is None:
            print("Użyj: config set <klucz> <wartość>")
            return 2
        value: Any = args.value
        if args.key in ("max_cpu_percent", "max_gpu_percent", "idle_minutes", "max_disk_gb", "max_task_seconds", "max_data_mb"):
            value = int(args.value)
        elif args.key in ("modalities", "pack_tiers"):
            value = [x.strip().upper() if args.key == "pack_tiers" else x.strip() for x in args.value.split(",") if x.strip()]
        elif args.key in ("leaderboard_opt_in",):
            value = args.value.lower() in ("1", "true", "yes")
        cfg[args.key] = value
        try:
            path = save(cfg)
        except ConfigError as exc:
            print(f"Nieprawidłowa wartość: {exc}")
            return 1
        print(f"Zapisano {args.key} w {path}.")
        return 0
    print("Użyj: config show | config set <klucz> <wartość>")
    return 2


def cmd_leaderboard(args) -> int:
    cfg = _load(args)
    try:
        data = _client(cfg).leaderboard(args.limit)
    except ApiError as exc:
        print(f"Błąd: {exc.message}")
        return 1
    if args.json:
        _print(data)
        return 0
    for item in data.get("items") or []:
        print(f"{item.get('rank'):>2}. {item.get('display_name')} · zadania {item.get('total_tasks')} · "
              f"aktualizacje {item.get('accepted_updates')} · {', '.join(item.get('badges') or [])}")
    if not data.get("items"):
        print("Brak wpisów. Ranking jest dobrowolny.")
    return 0


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(prog="nieboguard-compute", description="NieboGuard Compute — publiczny klient rozproszonego treningu AI.")
    parser.add_argument("--version", action="version", version=f"%(prog)s {__version__}")
    parser.add_argument("--config", help="ścieżka do pliku konfiguracji")
    parser.add_argument("--data-dir", help="katalog danych klienta")
    sub = parser.add_subparsers(dest="command", required=True)

    p = sub.add_parser("doctor", help="sprawdź system, sprzęt i backend")
    p.add_argument("--json", action="store_true")
    p.set_defaults(func=cmd_doctor)

    p = sub.add_parser("enroll", help="zarejestruj urządzenie po jawnej zgodzie")
    p.add_argument("--invite", required=True)
    p.add_argument("--consent", action="store_true", help="potwierdza przeczytanie i akceptację zgody")
    p.add_argument("--modality", choices=["audio", "vision", "both"], default="audio")
    p.add_argument("--max-cpu", type=int, default=50)
    p.add_argument("--max-gpu", type=int, default=60)
    p.add_argument("--schedule", choices=["idle", "always", "manual"], default="idle")
    p.add_argument("--idle-minutes", type=int, default=10)
    p.add_argument("--max-data-mb", type=int, default=200, help="maks. rozmiar paczki danych w MB (0 = bez paczek)")
    p.add_argument("--pack-tiers", default="MINI,STANDARD", help="dozwolone klasy paczek: MINI,STANDARD,EXTENDED,GPU")
    p.add_argument("--name", help="dobrowolna nazwa do rankingu")
    p.add_argument("--leaderboard", action="store_true")
    p.set_defaults(func=cmd_enroll)

    p = sub.add_parser("import-config", help="zaimportuj konfigurację pobraną ze strony")
    p.add_argument("file")
    p.set_defaults(func=cmd_import_config)

    p = sub.add_parser("agent", help="lokalny agent dla przycisków na stronie /compute/")
    p.add_argument("--host", default="127.0.0.1")
    p.add_argument("--port", type=int, default=8765)
    p.set_defaults(func=cmd_agent)

    p = sub.add_parser("start", help="rozpocznij pracę (wymaga zgody)")
    p.add_argument("--once", action="store_true", help="wykonaj najwyżej jedno zadanie")
    p.add_argument("--max-tasks", type=int, default=0)
    p.add_argument("--dev-allow-unsigned", action="store_true", help="tylko test lokalny; bez podpisu manifestu")
    p.set_defaults(func=cmd_start)

    p = sub.add_parser("stop", help="poproś o zatrzymanie po checkpointcie")
    p.set_defaults(func=cmd_stop)
    p = sub.add_parser("resume", help="usuń żądanie zatrzymania")
    p.set_defaults(func=cmd_resume)
    p = sub.add_parser("status", help="pokaż wkład i zadania")
    p.add_argument("--json", action="store_true")
    p.set_defaults(func=cmd_status)
    p = sub.add_parser("revoke", help="cofnij dostęp urządzenia")
    p.add_argument("--reason", default="volunteer requested revocation")
    p.set_defaults(func=cmd_revoke)
    p = sub.add_parser("config", help="pokaż lub zmień konfigurację")
    p.add_argument("action", choices=["show", "set"])
    p.add_argument("key", nargs="?")
    p.add_argument("value", nargs="?")
    p.set_defaults(func=cmd_config)
    p = sub.add_parser("leaderboard", help="dobrowolny ranking")
    p.add_argument("--limit", type=int, default=50)
    p.add_argument("--json", action="store_true")
    p.set_defaults(func=cmd_leaderboard)
    return parser


def main(argv: list[str] | None = None) -> int:
    parser = build_parser()
    args = parser.parse_args(argv)
    try:
        return int(args.func(args))
    except (ConfigError, RunnerError) as exc:
        print(f"Błąd: {exc}")
        return 1
    except BrokenPipeError:
        return 0


if __name__ == "__main__":
    raise SystemExit(main())
