# Publiczne API

Baza: `https://api.nieboguard.pl`

Publiczne endpointy są **read-only** i nie wymagają tokenu. Nie wysyłaj do nich
danych osobowych, sekretów ani współrzędnych prywatnych instalacji.

## Endpointy

| Metoda | Ścieżka | Opis |
|---|---|---|
| GET | `/api/state` | bieżący snapshot mapy; 20-sekundowy cache brzegowy |
| GET | `/api/v2/feed?limit=20` | komunikaty i zdarzenia |
| GET | `/api/v2/alarms?limit=20` | kanoniczne alarmy/ostrzeżenia |
| GET | `/api/v2/incidents?limit=20` | opublikowane incydenty |
| GET | `/api/v2/stats/overview` | statystyki rejestru |
| GET | `/api/v2/sources` | źródła, licencje i status |
| GET | `/api/v2/sources/:id` | szczegóły źródła |
| GET | `/api/v2/locations` | województwa z danymi |
| GET | `/api/v2/threats` | taksonomia zagrożeń |

## Przykład

```bash
curl -s https://api.nieboguard.pl/api/v2/stats/overview | jq .
curl -s 'https://api.nieboguard.pl/api/v2/feed?limit=5' | jq '.items[].title'
```

## Odpowiedzi i błędy

- Sukces: JSON z polami zależnymi od endpointu.
- Błąd: `{ "error": { "code": "...", "message": "..." } }`.
- Kody HTTP: 200, 204, 400, 404, 429, 503.
- Nie zakładaj stałej liczby rekordów ani stałej dostępności źródła.
- Nagłówki `Cache-Control` mogą się zmieniać; stosuj backoff i nie omijaj
  limitów.

## Rate limiting

Publiczne API może ograniczać częstotliwość. Domyślnie odpytuj `/api/state`
nie częściej niż co 20–30 s, a feed/statystyki co kilka minut. Przy `429`
zastosuj wykładniczy backoff.

## Bezpieczeństwo i etyka

- Nie publikuj danych osobowych, prywatnych lokalizacji ani nagrań bez zgody.
- Nie przedstawiaj odpowiedzi API jako oficjalnego alarmu.
- Nie usuwaj atrybucji źródeł i licencji.
- Nie próbuj omijać limitów, autoryzacji ani paneli administracyjnych.

## OpenAPI

Minimalny kontrakt: [`openapi.yaml`](openapi.yaml). To dokumentacja, nie
gwarancja niezmienności pól.
