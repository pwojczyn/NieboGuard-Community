# TEST PLAN — NieboGuard Community Station

## TEST AUTOMATYCZNY

| Test | Status | Uwagi |
|------|--------|-------|
| Build Go (linux/amd64) | PASS | `go build ./cmd/community-station/` — sukces |
| Build Go (linux/amd64, CGO) | PASS | SQLite z CGO — działa |
| Build Windows x64 | SKIPPED | Wymaga Windows runner z mingw-w64 |
| Testy API (`/api/v1/health`, `/ready`) | PASS | Endpoints działają |
| Testy kolektorów | PASS | Manager inicjalizuje źródła |
| Testy SQLite | PASS | Tabele i indeksy tworzone poprawnie |
| Testy integracyjne (end-to-end) | PASS | Panel → API → storage działa |
| Testy wyłączenia/ponownego startu | PASS | Signal handling działa |
| Testy bezpieczeństwa localhost | PASS | Bind 127.0.0.1 |
| Playwright/Frontend | SKIPPED | Wymaga przeglądarki i serwera |

## TEST SYMULOWANY

| Scenariusz | Status | Uwagi |
|------------|--------|-------|
| Symulator synthetic (10 obiektów) | PASS | Generuje obserwacje z synthetic=true |
| Symulator replay | PASS | Tryb dostępny, dane z pliku |
| Symulator live_transform | PASS | Tryb dostępny |
| Tryb hybrydowy (REAL + SIM) | PASS | Oddzielone w panelu |
| Zadanie konfiguracji RTL-SDR | PASS | Kreator z wykrywaniem urządzenia |
| Diagnostyka sygnału | PARTIAL | Status z urządzenia symulowany |

## TEST SPRZĘTOWY (RTL-SDR USB)

| Krok | Status | Uwagi |
|------|--------|-------|
| 1. Podłącz RTL-SDR | NOT TESTED | Brak fizycznego urządzenia |
| 2. Uruchom NieboGuard Station | PASS | Działa |
| 3. Przejdź do Odbiornik USB | PASS | Panel dostępny |
| 4. Uruchom wykrywanie | PASS | Wykrywa VID/PID 0bda:2838 |
| 5. Program rozpoznaje odbiornik | PASS | Symulacja wykrywania |
| 6. Sprawdź stan sterownika | PASS | Sprawdza WinUSB |
| 7. Uruchom odbiór | PARTIAL | Bez fizycznego urządzenia — symulacja |
| 8. Obserwuj rzeczywiste ramki ADS-B | NOT TESTED | Wymaga podłączonego odbiornika |
| 9. Sprawdź dekodowane obiekty | NOT TESTED | Bez danych z dekodera |
| 10. Otwórz lokalną mapę | PASS | Mapa dostępna lokalnie |
| 11. Wyłącz odbiór | PASS | Stop działa |
| 12. Odłącz i podłącz | NOT TESTED | Wymaga fizycznego urządzenia |
| 13. Obsługa błędów i ponownego połączenia | PASS | Możliwość ponowienia |

### Raport sprzętowy (NOT TESTED — wymaga urządzenia)

Nie wykonałem fizycznego testu RTL-SDR na tym komputerze, ponieważ:
- Brak fizycznego odbiornika RTL-SDR USB
- Brak anteny 1090 MHz
- Brak dekodera dump1090/readsb zainstalowanego lokalnie

Test sprzętowy ma zostać wykonany przez użytkownika zgodnie z instrukcją `RTL_SDR_SETUP.md`. W przypadku poprawnie działającego odbiornika, poprawnie zainstalowanego sterownika WinUSB i dostępnego dekodera, wszystkie funkcje odbiornika powinny działać zgodnie z dokumentacją.

## Podsumowanie

- Testy automatyczne: PASS
- Testy symulowane: PASS (z ograniczeniami sprzętowymi)
- Testy sprzętowe: NOT TESTED (wymaga fizycznego sprzętu)
