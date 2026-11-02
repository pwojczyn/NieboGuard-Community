# USB Troubleshooting — NieboGuard Community Station

## Błąd: "Device not found" / urządzenie niewidoczne

### Przyczyny i rozwiązania:

1. **Łączność USB**
   - Sprawdź czy kabel USB działa (przetestuj w innym portu, innym komputerze)
   - Porty USB 3.0 mogą mieć problemy z niekompatybilnymi urządzeniami — spróbuj USB 2.0
   - Spróbuj innego kabla (najlepiej krótki)

2. **Sterownik nie zainstalowany**
   - Panel → Odbiornik USB → Sprawdź sterownik
   - Jeśli WinUSB: zainstaluj przez Zadig (pat. RTL_SDR_SETUP.md)
   - Jeśli wyświetla się "Unknown Device" — aktualizuj sterowniki ręcznie

3. **Inny program używa urządzenia**
   - Zamknij inne aplikacje z RTLSDR (SDR#, SDR-Trunk, VirtualRadar)
   - Sprawdź w Task Managerze czy `rtl-sdr` nie działa w tle (choć aplikacja używa własnych adapterów)
   - Endpoint może być zajęty przez system — restart pomaga

4. **Wielokrotne urządzenia**
   - Jeśli masz więcej niż jedno RTL-SDR — w panelu wybierz odpowiedni przez numer urządzenia
   - VID/PID w info o urządzeniu pomoże odróżnić

### Rozwiązanie: Reset urządzenia
1. Odłącz RTL-SDR z USB
2. Odczekaj 10 sekund
3. Podłącz ponownie
4. Otwórz panel → Odbiornik USB → Wybierz ponownie
5. Kliknij "Test urządzenia"

## Błąd: "No aircraft messages" / "0 msg/s"

### Przyczyny:
- Brak anteny lub słaba jakość sygnału
- Gain ustawiony błędnie (za niskie = nie słyszymy, za wysokie = nasycenie)
- Zbyt duża odległość od lotniska / samolotów
- Zakłócenia od innych urządzeń RF

### Rozwiązania:
1. Sprawdź gain w panelu (auto to dobry punkt startowy)
2. Przesuń antenę na zewnątrz, na dach (najlepsze wyniki)
3. Wyłącz inne urządzenia generujące zakłócenia (lampy energoświetlne, ładowarki USB)
4. Sprawdź PPM correction (korekcja często poprawia wyniki w niektórych lokalizacjach)

## Błąd: "PPM correction" — korekcja częstotliwości

### Co to jest:
PPM (Parts Per Million) to błąd częstotliwości odbiornika.
Domyślne 0 często działa, ale niektóre odbiorniki mają bias ~20-50 PPM.

### Jak ustawić:
- Wartości dodatnie (+) przesuwają częstotliwość w górę
- Wartości ujemne (-) przesuwają w dół
- Zmieniaj krokami po 5 PPM
- Obserwuj liczbę odebranych ramek — im więcej, tym lepiej

## Błąd: przeglądarka nie otwiera panelu

### Rozwiązanie:
1. Uruchom ręcznie: `http://127.0.0.1:8787`
2. Sprawdź czy port 8787 nie jest blokowany przez firewall
3. Uruchom z flagą: `.\NieboGuard-Station.exe --port 8788` (użyj innego portu)
4. Zablokuj dostęp z zewnątrz w ustawieniach aplikacji

## Diagnostyka systemowa

### Sprawdź logi aplikacji
- Panel → Diagnostyka → Logi
- Poszukiwaj wpisów z czerwonymi znakami ERROR

### Sprawdź porty
```powershell
netstat -an | findstr 8787
```
Jeśli nic nie wyjdzie — port jest wolny.

### Sprawdź procesy
- Panel → Ustawienia → Diagnostyka → Status procesów
- Aplikacja powinna uruchomić nnode.exe wewnętrznie

## Kontakt z pomocą

Jeśli żadne rozwiązanie nie pomoże:
1. Wygeneruj pakiet diagnostyczny: Panel → Diagnostyka → Eksport pakietu
2. Wyślij together z logami do grupy wsparcia
3. W opisie problemu podaj:
   - Wersję aplikacji (Footer panelu)
   - Model RTL-SDR (np. R820T2)
   - Sterownik (WinUSB? Other?)
   - Wynik testu "Sprawdź sterownik"