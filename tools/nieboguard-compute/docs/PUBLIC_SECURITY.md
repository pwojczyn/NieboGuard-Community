# NieboGuard Compute — bezpieczeństwo i model zagrożeń

Data: 2026-11-01. Dokument opisuje stan MVP, w tym ograniczenia.

## Zasady

1. Backend nie dostarcza kodu do wykonania.
2. Manifest zadania jest podpisany Ed25519; klient ma przypięty klucz
   publiczny. Brak klucza = brak startu, chyba że jawny tryb testowy.
3. Wynik nigdy nie trafia bezpośrednio do produkcji. Najpierw kwarantanna,
   potem walidacja i niezależna ewaluacja.
4. Publikacja i rollback są ręczne.
5. Nie ufamy metrykom, nazwom tensorów, typom ani rozmiarom deklarowanym przez
   klienta — serwer sprawdza artefakt.
6. Minimalizacja danych: system, architektura, rdzenie, kubełek RAM, model i
   kubełek VRAM. Bez nazwy hosta, MAC, serii, ścieżek, GPS i prywatnych mediów.
7. Nie deklarujemy pełnej anonimowości aktualizacji wag. Urządzenie ma
   identyfikator i historię zadań.

## Uwierzytelnianie urządzeń

- Rejestracja wymaga zaproszenia (hash w D1), wersji zgody i jawnego
  `accept_consent: true`.
- Refresh token: 32 losowe bajty, przechowywany jako SHA-256. Pokazywany
  wyłącznie raz.
- Access token: 32 losowe bajty, TTL 20 min (konfigurowalne 5–60), scopes:
  `task:read`, `update:write`, `stats:read`, `consent:revoke`.
- Cofnięcie urządzenia unieważnia refresh i wszystkie access tokeny oraz
  anuluje aktywne leasingi.
- Rate limiting: rejestracja, odświeżanie, pobieranie zadania, zgłoszenia
  wyników, uploady.
- Brak tokenów w URL poza capability URL do pojedynczego obiektu R2.

## Capability URL do R2

Projekt **nie ma** bindingu R2 S3/presigned URL. Zamiast tego Worker wystawia
krótkotrwały URL:

```text
/api/v2/compute/objects/<upload_id>?token=<payload>.<HMAC>
```

Token zawiera: upload_id, purpose, device_id, task_id, direction, exp,
max_bytes, content_type, expected SHA-256 i nonce. Podpis HMAC-SHA-256 jest
weryfikowany przed odczytem/zapisem. Slot uploadu jest jednorazowy
(`status=PENDING` → `COMPLETED`). Serwer egzekwuje:

- metodę PUT/GET i kierunek,
- rozmiar z `Content-Length` i twardy limit,
- typ MIME,
- SHA-256 ciała i oczekiwaną sumę,
- status urządzenia (aktywne/cofnięte),
- zgodność zadania i rundy.

**Ograniczenie:** to nie jest presigned S3. Cały ruch uploadu przechodzi przez
Worker, więc limit rozmiaru i czasu pracy Workera jest jednocześnie limitem
MVP. Dla większych modeli trzeba dodać R2 S3 API i presigned URL po audycie.

## Walidacja safetensors

- tylko `F32` w kliencie; serwer statystyk obsługuje też F16/BF16,
- limit 16 MB i 1 MB nagłówka,
- brak pickle, ZIP, wykonywalnych plików i kodu,
- dokładna zgodność nazw, dtype i kształtów z manifestem,
- brak NaN/Inf,
- limit `abs_max` i dryfu L2 względem statystyk bazy,
- wymagane metadane: task, round, model bazowy, wersja, hash architektury,
  klasy, liczba próbek,
- odrzucenie duplikatów SHA-256.

## Ochrona przed poisoningiem

- zadanie i runda muszą się zgadzać,
- architektura i model bazowy muszą być identyczne,
- zakresy wag i dryf od bazy,
- klasy i liczba próbek w metadanych,
- niezależny holdout i raport evaluatora,
- progi per klasa i fałszywe alarmy,
- trimmed mean zamiast zwykłej średniej,
- ręczny przegląd przypadków granicznych,
- publikacja tylko po jawnej decyzji administratora i możliwość rollbacku.

To nie gwarantuje odporności na zaawansowany atak. Publiczna beta wymaga
zewnętrznego audytu.

## Prywatność

- MVP nie wysyła prywatnych zdjęć ani nagrań.
- Datasety wolontariuszy muszą mieć `redistribution_allowed=1` i
  `training_allowed=1`; inaczej Worker odmawia utworzenia zadania.
- Klient sprawdza prawa w manifescie po raz drugi.
- Rejestr danych zawiera autorów, licencję, pochodzenie, wersję i SHA-256.
- Wariant federacyjny z prywatnymi danymi jest osobnym, przyszłym etapem.
- Aktualizacje wag mogą ujawniać informacje o danych treningowych; nie
  obiecujemy anonimowości bez dodatkowych mechanizmów (DP, secure aggregation).

## Sekrety

Klucz prywatny manifestu, sekret uploadu, token evaluatora i konfiguracja
Workera pozostają wyłącznie w prywatnym backendzie. Publiczny klient otrzymuje
tylko klucz publiczny i krótkotrwałe tokeny urządzenia. Nie umieszczaj tych
wartości w repozytorium, issue, PR, logach ani zrzutach ekranu.

## Incydenty

- Cofnij urządzenie w panelu (`/admin/compute/devices/:id/revoke`).
- Unieważnij zaproszenia.
- Zablokuj zadanie/rundę.
- Oznacz model `REVOKED` i wykonaj rollback.
- Zachowaj `compute_audit`, artefakty i raporty do analizy.
- Nie usuwaj dowodów przed analizą.

## Znane ograniczenia

- Brak zewnętrznego audytu bezpieczeństwa.
- Brak podpisanych instalatorów.
- Brak testu na fizycznym GPU NVIDIA.
- Brak secure aggregation / differential privacy.
- Worker nie uruchamia niezależnej ewaluacji; wymaga zaufanego runnera.
- Panel admina korzysta z istniejącej sesji; brak osobnych ról dla Compute.
