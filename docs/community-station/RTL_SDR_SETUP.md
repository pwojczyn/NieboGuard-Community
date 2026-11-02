# RTL-SDR Setup — NieboGuard Community Station

## Wymagania sprzętowe
- RTL2832U / RTL2832U (R820T2)
- Antena 1090 MHz (np. dipol 1/4 fali)
- Kabel USB 2.0 (nie za długi — <3m preferowane)

## Krok 1: Podłącz odbiornik
- Wklej RTL-SDR do portu USB
- Sprawdź czy Windows widzi urządzenie (Menedżer Urządzeń)
- Jeśli widzisz błędy yellow — znaczy to, że sterownik jest nieprawidłowy

## Krok 2: Sprawdź sterownik
Panel → Odbiornik USB → "Sprawdź sterownik"
- Jeśli WinUSB = false → zainstaluj przez Zadig (pat. RTL_SDR_SETUP.md)
- Jeśli widzisz "Unknown Device" — znaczy to, że sterownik jest nieprawidłowy

## Krok 3: Zadig (jeśli potrzebne)
1. Pobierz: https://zadig.akeo.ie/
2. Uruchom jako Administrator
3. Options → List All Devices
4. Wybierz urządzenie z VID=0bda, PID=2838 (np. "RTL2832U")
5. W polu "Driver" wybierz WinUSB
6. Kliknij Replace Driver
7. Restart aplikacji

⚠️ Uwaga: Nie zmieniaj sterownika innych urządzeń USB!

## Krok 3: Konfiguracja odbiornika
W panelu → Odbiornik USB:
- Wybierz urządzenie (jeśli masz kilka — wskaż przez numer)
- Przycisk "Wykryj odbiornik" — pokaże nazwę i VID/PID
- Przycisk "Test urządzenia" — sprawdza połączenie
- Częstotliwość: 1090000000 Hz (1090 MHz)
- Próbka: 2400000 Hz (2.4 MHz)
- Gain: Auto (lub manual — np. 30-45 dB, eksperymentuj)
- PPM correction: 0 (dla większości odbiorników)

## Krok 5: Start odbioru
- Kliknij "Start odbioru"
- Status zmienia się na RECEIVING
- Liczba wiadomości/s rośnie
- Obiekty pojawiają się w sekcji Map

## Diagnostyka problemów

### Brak wiadomości
- Sprawdź antenę i jej pozycję (najlepiej zewnętrzna, powyżej dachu)
- Sprawdź gain w panelu (auto to dobry punkt startowy)
- Sprawdź PPM correction (0 to domyślne, ale możesz spróbować ±10)
- Upewnij się, że 1090 MHz to właściwa częstotliwość (ADS-B)

### Błędy CRC / mało ramek
- To normalne — nie wszystkie wiadomości są poprawne
- Jeśli błędy dominują — sprawdź antenę, zbliż się do lotniska
- Sprawdź PPM correction — korekcja często poprawia dekodowanie

### Urządzenie zajęte przez inny program
- Zamknij inne aplikacje RAD (e.g. SDR#, HDSDR, CubicSDR)
- Sprawdź w Task Managerze czy `dump1090.exe` działa w tle
- Endpoint może być zajęty przez system — restart aplikacji pomaga

## Tryb offline / symulator
Jeśli nie masz odbiornika — włącz symulator:
- Ustawienia → Symulator → Enabled = true
- Wybierz scenariusz: synthetic/replay/live_transform
- Obiekty SIM pojawią się na mapie z oznaczeniem SIM
- Dane REAL i SIM są wyraźnie rozdzielone (SIM = symulowane=true)
- W panelu widać różnicę: aktywne (REAL) vs symulowane (SIM)