# NieboGuard Community Station — Windows Install

## Wymagania systemowe
- Windows 10/11 x64
- Uprawnienia administratora (dla instalacji sterownika WinUSB)
- RTL-SDR USB (opcjonalnie, do odbioru rzeczywistego)
- Antena 1090 MHz (opcjonalnie)

## Pobieranie i uruchamianie

### 1. Pobranie
```powershell
# Z GitHub Releases
# https://github.com/pwojczyn/NieboGuard-Community/releases
```

### 2. Uruchomienie
```powershell
# Dwuklik na NieboGuard-Station.exe
# Lub z linii poleceń:
.\NieboGuard-Station.exe

# Z własną konfiguracją:
.\NieboGuard-Station.exe --config C:\Users\USER\NieboGuard\Station\config.yaml

# Bez otwierania przeglądarki:
.\NieboGuard-Station.exe --nopanel

# Na innym porcie:
.\NieboGuard-Station.exe --port 8788
```

### 3. Pierwsze uruchomienie
1. Aplikacja tworzy katalog `%LOCALAPPDATA%\NieboGuard\Station`
2. Tworzy domyślny `config.yaml`
3. Inicjalizuje bazę SQLite `community_station.db`
4. Uruchamia serwer HTTP na porcie 8787
5. Otwiera panel w przeglądarce: `http://127.0.0.1:8787`

## Konfiguracja RTL-SDR

### 1. Podłącz odbiornik USB
- Wklej RTL-SDR do portu USB
- Windows może zainstalować domyślny sterownik (często nieprawidłowy)

### 2. Sprawdź sterownik WinUSB
W panelu → Ustawienia → Odbiornik USB → "Sprawdź sterownik"
- Jeśli "WinUSB installed: false" → konieczna instalacja Zadig

### 3. Instalacja Zadig (jeśli wymagane)
1. Pobierz Zadig: https://zadig.akeo.ie/
2. Uruchom jako Administrator
3. Options → List All Devices
4. Wybierz "RTL2832U" (VID: 0bda, PID: 2838)
5. Wybierz WinUSB jako sterownik
6. Kliknij "Replace Driver"
7. Restart aplikacji

### 4. Konfiguracja odbioru
W panelu → Odbiornik USB:
- "Wykryj odbiornik" → wykrywa urządzenie
- "Test urządzenia" → sprawdza połączenie
- Ustaw Gain: Auto / Manual
- PPM correction (opcjonalnie)
- Kliknij "Start odbioru"

## Rozwiązywanie problemów

### Odbiornik nie wykrywany
- Sprawdź połączenie USB (inny port/kabel)
- Sprawdź Menedżer urządzeń → "libusb-win32 devices"
- Upewnij się, że inny program nie używa urządzenia (dump1090, SDR#, itp.)

### Brak ram ADS-B
- Sprawdź antenę (1090 MHz,'extérieur preferred)
- Sprawdź gain (auto lub manual ~30-45 dB)
- Sprawdź PPM (korekcja częstotliwości)
- Brak samolotów w zasięgu = normalne w niektórych lokalizacjach

### Błąd "Device busy"
- Zamknij inne programy używające RTL-SDR
- Restart aplikacji
- Odłącz i podłącz USB

### Panel nie otwiera się
- Sprawdź czy port 8787 jest wolny: `netstat -an | findstr 8787`
- Użyj `--port 8788`
- Sprawdź firewall (domyślnie tylko localhost)

### Bazę danych nie można otworzyć
- Sprawdź uprawnienia do `%LOCALAPPDATA%\NieboGuard\Station`
- Usuń `community_station.db` i uruchom ponownie (utracisz historię)

## Aktualizacja

1. Pobierz nową wersję EXE
2. Zastąp stary plik
3. Uruchom - baza i konfiguracja są zachowane

## Odinstalowanie

1. Zamknij aplikację
2. Usuń plik EXE
3. Opcjonalnie usuń `%LOCALAPPDATA%\NieboGuard\Station`

## Zależności zewnętrzne (opcjonalnie)

- **dump1090/readsb**: Jeśli chcesz używać zewnętrznego dekodera zamiast wbudowanego
  - Zainstaluj w `C:\Program Files\dump1090\`
  - Wskazuj ścieżkę w konfiguracji
- **Zadig**: Do instalacji sterownika WinUSB (tylko raz)

## Dokumentacja dodatkowa

- RTL_SDR_SETUP.md — szczegółowa instrukcja
- USB_TROUBLESHOOTING.md — typowe problemy
- ARCHITECTURE.md — architektura aplikacji
- API.md — REST API i WebSocket