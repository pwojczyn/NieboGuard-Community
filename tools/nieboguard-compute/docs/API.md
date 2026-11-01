# NieboGuard Compute — publiczny kontrakt klienta

Bazowy prefiks: `/api/v2`. Backend jest prywatny; ten dokument opisuje wyłącznie
część potrzebną klientowi. Wszystkie błędy mają postać
`{"error":{"code":"...","message":"..."}}`.

## Publiczne

### `GET /compute/status`

`enabled`, `configured`, `manifest_public_key`, `manifest_key_id`, liczba
uczestników, aktywni 24 h, podział OS, GPU NVIDIA, statusy zadań/rund, liczba
modeli i datasetów, wersja zgody. Bez dokładnego sprzętu pojedynczych urządzeń.

### `GET /compute/leaderboard?limit=50`

Tylko urządzenia, które wyraziły zgodę i podały nazwę. Bez ID i sprzętu.

### `POST /compute/enroll`

```json
{
  "invite_code": "ngi_...",
  "accept_consent": true,
  "consent_version": "2026-11-01",
  "hardware": {
    "os_family": "linux", "arch": "x64", "cpu_cores": 8, "ram_gb": 16,
    "gpu_vendor": "nvidia", "gpu_model": "RTX 3060", "gpu_memory_gb": 8
  },
  "preferences": {
    "modalities": ["audio", "vision"],
    "max_cpu_percent": 50,
    "max_gpu_percent": 60,
    "schedule": "idle",
    "max_data_mb": 200,
    "pack_tiers": ["MINI", "STANDARD"]
  },
  "display_name": "opcjonalnie",
  "leaderboard_opt_in": false,
  "capabilities": {"backends": ["python"]},
  "client": {"version": "0.1.0", "protocol": 1}
}
```

Odpowiedź 201: `device_id`, jednorazowy `refresh_token`, `access_token`,
`access_expires_at`, `scopes`, `manifest_public_key`, `manifest_key_id`.
Błędy: zgoda/walidacja 400, zaproszenie 403, wyczerpane 409, limit 429,
pojemność 503.

### `POST /compute/token`

`{device_id, refresh_token}` → krótkotrwały `access_token`.

### `GET /compute/me` (Bearer, `stats:read`)

Statystyki wkładu i ostatnie zadania.

### `PUT /compute/me/preferences` (Bearer, `preferences:write`)

Aktualizuje modalności, limity CPU/GPU, harmonogram, `max_data_mb` i
`pack_tiers` po rejestracji.

### `POST /compute/consent/revoke` (Bearer, `consent:revoke`)

Unieważnia urządzenie, refresh i access tokeny oraz anuluje leasingi.

## Zadania

### `GET /compute/tasks/next` (Bearer, `task:read`)

`{ok:true, task:null}` albo `{ok:true, resumed, manifest}`.

Manifest jest podpisany Ed25519. Zawiera:

- `task_id`, `round_id`, `assignment_id`, termin leasingu i limit czasu,
- `dataset`: ID, wersja, licencja, atrybucja, prawa, SHA-256, shardy,
- opcjonalny `pack`: ID, klasa MINI/STANDARD/EXTENDED/GPU, bajty, SHA-256,
  liczba sesji, osobne `model_bytes` i `library_bytes`,
- `base_model`: dokładna wersja, SHA-256, hash architektury, spec tensorów,
- `classes`, `hyperparams`, `limits`,
- `upload`: jednorazowy upload_id, URL, typ, limit,
- `signature`: `{alg:"Ed25519", key_id, value}`.

Każdy shard ma własny URL capability, rozmiar i SHA-256. Klient sprawdza
origin, rozmiar, sumę i limit paczki przed treningiem.

### `POST /compute/tasks/:id/complete`

```json
{
  "upload_id": "cmu_...",
  "sha256": "...",
  "bytes": 1509,
  "metrics": {
    "samples": 100,
    "duration_s": 31.2,
    "accuracy": 0.92,
    "loss": 0.12,
    "epochs_completed": 6,
    "backend": "python"
  }
}
```

Wynik `202` = przyjęty do kwarantanny i oczekuje na niezależną ewaluację.
`422` = odrzucony przed ewaluacją. Metryki klienta nie są dowodem jakości.

### `POST /compute/tasks/:id/fail`

`{error}` — zgłoszenie niepowodzenia.

## Obiekty

`GET/PUT /compute/objects/:uploadId?token=...` — krótkotrwały capability URL
do pojedynczego obiektu. Wymagany rozmiar, typ i SHA-256. Klient nie używa
adresów spoza skonfigurowanego origin API.

## Ocena i publikacja

Niezależna ewaluacja, agregacja, publikacja i rollback są po stronie
prywatnego backendu i administratora. Klient nie ma endpointu publikacji i nie
może zatwierdzić własnego wyniku.
