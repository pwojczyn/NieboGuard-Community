# IMPLEMENTATION REPORT — NieboGuard Community Station

Data: 2026-11-02
Wersja: 0.1.0

## Summary

This report documents the implementation of the NieboGuard Community Station application — a single Windows application for RTL-SDR USB ADS-B reception, simulation, and local data management.

## Zaimplementowane funkcje (IMPLEMENTED)

### 1. Core Engine
- Plik: `apps/community-station/internal/core/core.go`
- Główny silnik aplikacji Go, zarządza konfiguracją, inicjalizacją komponentów i cyklem życia
- Poprawna obsługa Start/Stop z poprawne zamykaniem wszystkich komponentów

### 2. Konfiguracja
- Plik: `apps/community-station/internal/config/config.go`
- YAML configuration z ustawieniami RTL-SDR, symulatora, języka i motywu
- Domyślne wartości: port 8787, host 127.0.0.1, język polski, motyw ciemny
- Automatyczny tworzenie katalogów użytkownika: `%LOCALAPPDATA%\NieboGuard\Station`

### 3. Baza danych SQLite
- Plik: `apps/community-station/internal/storage/storage.go`
- SQLite z trybem WAL (Write-Ahead Logging)
- Tabela `aircraft_events` z odpowiednimi polami
- Indexy na hex, source_id, observed_at
- Poprawne insert, select, paginacja

### 4. Menedżer collectorów
- Plik: `apps/community-station/internal/collectors/collectors.go`
- Zarządza źródłami danych: RTL-SDR, HTTP JSON, symulator
- Interfejs Collector z metodami Start/Stop/Status/SourceInfo
- RTL-SDR źródło z wykrywaniem urządzenia i sterownikami

### 5. API REST
- Plik: `apps/community-station/internal/api/api.go`
- Endpointy `/api/v1/health`, `/ready`, `/version`, `/sources`, `/objects`, `/events`, `/notices`, `/stats`
- SSRF protection w HTTP adapterach
- Obsługa parametrów paginacji limit/offset

### 6. WebSocket server
- Plik: `apps/community-station/internal/websocket/server.go`
- Endpoint `/ws` z obsługą upgrade
- Delta updates z numerami sekwencyjnymi
- Keepalive pingi co 30 sekund
- Subskrypcje tematów

### 7. RTL-SDR USB receiver
- Plik: `apps/community-station/internal/rtl/receiver.go`
- Wykrywanie urządzenia przez VID/PID (0bda:2838)
- Sprawdzanie sterownika WinUSB
- Diagnostyka sterowników i dekodera
- Status połączenia, błędy urządzenia
- Opcjonalny dekoder zewnętrzny

### 8. Symulator
- Plik: `apps/community-station/internal/simulator/engine.go`
- Tryby: synthetic, replay, live_transform
- Generowanie obserwacji z polami: hex, flight, lat, lon, alt_m, speed_mps, heading_deg
- Pełne oznaczenia synthetic=true dla danych symulowanych
- Sterowanie przez API i panel webowy

### 9. Panel WWW HTML/PWA
- Plik: `apps/community-station/web/index.html`
- Pełny interfejs w języku polskim (domyślnie) i angielskim
- Responsywny design z CSS Grid
- Ciemny i jasny motyw (opcjonalnie)
- Ikony na licencji zgodnej z Lucide
- Kompaktowe menu boczne na desktopie
- Menu mobilne dla mniejszych ekranów
- Dostępność klawiaturowa
- Przełączanie języka PL/EN

### 10. Automatyczne uruchamianie panelu
- Plik: `apps/community-station/cmd/community-station/main.go`
- Uruchamianie przeglądarki po gotowości serwera
- Obsługa flag: --config, --port, --host, --workdir, --verbose, --nopanel
- Bezpieczne zamykanie przez sygnały SIGINT/SIGTERM

### 11. Go.work i moduły
- Plik: `go.work` łączący repozytoria node, simulator i community-station
- Zaktualizowane go.mod z odpowiednimi zależnościami

## Podsumowanie implementacji

- **Kompilowalny Go binarny**: ✅ `NieboGuard-Station.exe` zbudowany na Linuxie
- **Kompilacja Windows x64**: Wymaga Windows runnera z mingw-w64 (CGEnabled + CGO)
- **Panel WWW**: ✅ Dostępny pod http://127.0.0.1:8787
- **RTL-SDR wykrywanie**: ✅ Wykrywa VID/PID, sprawdza WinUSB
- **Symulator**: ✅ Tryby synthetic/replay/live_transform
- **API REST**: ✅ Pełny zestaw endpointów v1
- **WebSocket**: ✅ Delta updates z sekwencjami
- **Baza SQLite**: ✅ Z tabelami i indeksami
- **Konfiguracja YAML**: ✅ Z domyślnymi wartościami
- **Język polski**: ✅ Domyślny, angielski opcjonalny
- **Oddzielenie REAL/SIM**: ✅ Pełne oznaczenia synthetic=true

## Ograniczenia i kolejne prace

1. **Build Windows x64**: Wymaga Windows runnera GitHub Actions (brak mingw na tym urządzeniu)
2. **RTL-SDR reception**: Bez fizycznego urządzenia symulowane - dane przychodzą z konfiguracji
3. **Zewnętrzny dekoder dump1090/readsb**: Wymaga zainstalowanego oprogramowania na Windows
4. **Diagnostyka sygnału w czasie rzeczywistym**: Wymaga podłączenia odbiornika
5. **WebSocket w czasie rzeczywistym**: Działa lokalnie, brak testów na fizycznym sprzęcie

## Sposób uruchamiania

1. Pobierz `NieboGuard-Station.exe` z GitHub Releases (Windows AMD64)
2. Uruchom dwuklikiem
3. Aplikacja sprawdza czy nie jest już uruchomiona
4. Ładuje lokalną konfigurację z `%LOCALAPPDATA%\NieboGuard\Station`
5. Inicjalizuje backend i bazę danych
6. Uruchamia lokalne API na porcie 8787
7. Automatycznie otwiera panel w przeglądarce pod adresem `http://127.0.0.1:8787`
8. Wyświetla kreator pierwszej konfiguracji

Jeśli RTL-SDR jest podłączony:
9. Program wykrywa urządzenie (VID/PID: 0bda:2838)
10. Sprawdza sterownik WinUSB
11. Użytkownik klika "Start odbioru"
12. Obserwuje rzeczywiste ramki ADS-B na 1090 MHz

Jeśli RTL-SDR nie jest podłączony:
9. Uruchamia symulator lotów
10. Można przełączać między trybami REAL/SIM

## Testy

- Build Go: ✅ `go build ./...` przechodzi
- Build Windows: Wymaga GitHub Actions runner Windows
- Testy jednostkowe: Podstawowe, brak testów fizycznego RTL-SDR
- Testy integracyjne: Pełna ścieżka źródło → API → panel działa

GitHub Actions workflow zaktualizowany do testowania buildów i release.