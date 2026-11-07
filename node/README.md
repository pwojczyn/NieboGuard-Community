# NieboGuard Community Node

Lokalny serwer danych dla NieboGuard. Domyślny adres: `127.0.0.1:8787`.

## Uruchomienie

```bash
cd node
go run ./cmd/nnode --host 127.0.0.1 --port 8787
```

## Bezpieczeństwo lokalnego API

Domyślnie dane i endpointy administracyjne są dostępne tylko z loopbacku.
Aby chronić je tokenem:

```bash
export NNODE_AUTH_TOKEN='<long-random-token>'
export NNODE_REQUIRE_AUTH=1
```

Klient wysyła:

```http
X-NG-Local-Auth: <token>
```

lub:

```http
Authorization: Bearer <token>
```

- `/api/v1/health`, `/api/v1/ready` i `/api/v1/version` są lokalnymi probe'ami.
- `/api/v1/objects`, `/sources`, `/events`, `/stats` wymagają tokenu, jeśli jest
  skonfigurowany.
- `/api/v1/admin/*` bez tokenu jest zawsze odrzucane.
- Limit: 120 żądań/minutę na adres; nagłówki `Retry-After` przy 429.
- Nie wystawiaj portu na Internet. Token nie jest zapisywany w repo.

Testy:

```bash
go test ./...
```

## Offline packs

Metadane paczek offline są przechowywane w SQLite:

```text
GET  /api/v1/offline/packs?id=pack-id
GET  /api/v1/offline/packs/{id}
POST /api/v1/admin/offline/packs
```

Admin wymaga `NNODE_AUTH_TOKEN`. API przyjmuje wyłącznie manifest i SHA-256;
nie przyjmuje ścieżki plików. Symulowana paczka musi mieć `LOCAL_ONLY`.
