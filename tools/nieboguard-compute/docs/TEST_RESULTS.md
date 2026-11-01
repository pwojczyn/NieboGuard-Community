# NieboGuard Compute — wyniki testów publicznego klienta

Data: 2026-11-01. Środowisko: Linux x86_64, Python 3.14.4, PowerShell 7.4.6
(tylko parser `.ps1`). Bez produkcji i bez prywatnego backendu.

## Testy jednostkowe klienta

```bash
cd tools/nieboguard-compute
python -m unittest discover -s tests -v
```

Wynik: **35 testów, OK**.

Zakres:

- RFC 8032 Ed25519 i canonical JSON zgodny z JS,
- safetensors roundtrip, NaN/Inf, trailing data, shape mismatch,
- manifest: podpis, wygaśnięcie, brak redystrybucji, obcy origin,
- paczki danych: limit `max_data_mb`, dozwolona klasa, SHA-256,
- konfiguracja: uprawnienia 0600, redakcja tokenu, HTTPS, kompletność danych,
- lokalny agent: status, stop/resume, odmowa startu bez rejestracji, blokada
  obcego Origin, preflight CORS, start/stop z zamockowanym procesem,
- Audio AI: WAV → cechy, trening CPU, checkpoint i wznowienie, split z manifestu,
  base artifact, stop bez uploadu,
- Vision AI: głowa detekcji CPU, objectness F1 ≥ 0.8, eksport safetensors,
- canonical JSON dla `1e-5`/`1e-6`/`1e-7`/`1e20`/`1e21`,
- przygotowanie paczek: ZIP/TAR, zip-slip, absolute path, grupy sesji, manifest.
- limit danych bez paczki i dozwolone klasy paczek.

## Instalator Linux

```bash
./install-linux.sh --skip-backend-check -y
```

Wynik: **OK** w czystej kopii. Utworzono `.venv`, wrapper
`bin/nieboguard-compute` i `data/config.json` z uprawnieniami właściciela.

## Skrypty Windows

Parser PowerShell 7.4.6:

```powershell
[System.Management.Automation.Language.Parser]::ParseFile(...)
```

Wynik: `install-windows.ps1`, `start-windows.ps1`, `uninstall-windows.ps1` —
**OK (składnia)**.

**Nie wykonano** testu runtime na Windows 10/11. To pozostaje otwarte.

## Czego nie przetestowano publicznie

- pełnego E2E z Workerem, D1 i R2 (test w prywatnym repo backendu),
- fizycznego GPU NVIDIA i CUDA,
- podpisanych instalatorów EXE/MSI/DEB/RPM,
- Playwright/UI E2E strony `/compute/`,
- rzeczywistych paczek danych z audytowanych źródeł.
