# INICJALNY AUDYT ARCHITEKTURY

Data: $(date)
Status: Complete

## 1. Badanie dostępnych repozytoriów

### 1.1 Główny prywatny repozytorium NieboGuard
- **Dostęp**: Prywatny, brak pełnego dostępu
- **Charakter**: Produkcyjny backend VPS, serwery HTTP/WSS, baza danych, mechanizmy synchronizacji
- **Ograniczenie**: Brak bezpośredniego dostępu do kodu produkcyjnego, konfiguracji VPS, kluczy API ani wewnętrznych tokena
- **Wniosek**: Musi zostać zaprojektowany niezależny interfejs integracyjny, który nie narusza granic prywatnego kodu, ale zachowuje kompatybilność na poziomie protokołów i formatów danych

### 1.2 Publiczne repozytorium Community (NieboGuard-Community)
- **Adres**: https://github.com/pwojczyn/NieboGuard-Community
- **Stan**: Dostępny pełen
- **Zawartość**:
  - Dokumentacja publicznego API i architektury
  - SDK JavaScript i Python (tylko do odczytu)
  - Narzędzia: ADS-B → GeoJSON, audyt audio, audyt obrazów
  - Syntetyczne dane testowe
  - Tutoriale dla wolontariuszy
  - Brak kodu produkcyjnego serwera
  - Brak binarek lub instalatorów
- **Wniosek**: To repozytorium pełni rolę zaplecza publicznego. Nowy projekt NieboGuard Node będzie rozwijany tutaj, ale jako całkowicie niezależny produkt open source.

### 1.3 Strona główna: https://nieboguard.pl/
- Cel: Oficjalna strona projektu, informacje produkcyjne
- Brak bezpośredniego połączenia z zadaniami implementacyjnymi

### 1.4 Projekt referencyjny: https://github.com/cukierrro/Straznik
- Jest to projekt "Straznik" (Strażnik)
- Architektura i kod wczesnych wersji (do 1.7.53) udostępniono na licencji MIT
- Od wersji 1.7.54 obowiązuje "inna, restrykcyjna licencja"
- **Kluczowa uwaga**: Nie wolno kopiować kodu z wersji objętych restrykcyjną licencją (po 1.7.54)
- Dostępne elementy do wykorzystania:
  - Koncepcja monitorowania przestrzeni powietrznej
  - Schematy formatów danych
  - Pomysły dotyczące adapterów źródeł
- Decyzja: Własna implementacja w Go, inspirowana koncepcjami, ale niezależna od restrykcyjnego kodu

## 2. Identyfikacja istniejących serwerów i mechanizmów (z głównego NieboGuard)

### 2.1 Serwery HTTP
- Dostępne endpointy read-only publicznego API
- Formaty: JSON, GeoJSON
- Statyczne odpowiedzmi bez autoryzacji do odczytu

### 2.2 Serwery WebSocket (WSS)
- Połączenia maintainence do centrali
- Subskrypcja zdarzeń na żywo
- Wymagana autoryzacja

### 2.3 Źródła danych (istniejące w projekcie referencyjnym)

| Źródło | Format | Dostępność | Uwagi |
|--------|--------|------------|-------|
| ADS-B (readsb/dump1090) | Beast TCP / JSON | Lokalny | Automatyczne wykrywanie portów |
| HTTP JSON | Aplikacja-specyficzny | Konfigurowalny | Mapowanie pól |
| RSS / Atom | XML | Publiczne | Komunikaty regionalne |
| Komunikaty RCB, RSO/TVP, PAŻP | Specjalne API | Ograniczone | Warunki prawne zweryfikowane |
| NEPTUN | Niejawny | Niepubliczne | Tylko interfejs adaptera |
| CelesTrak | Publiczne | Tak | Dane orbitalne |

### 2.4 Formaty danych (istniejące w repozytorium)
- `aircraft`: hex, flight, lat, lon, alt_baro, gs, track, seen
- `feed items`: id, kind, source, category, title, region, verification, published_at
- `sample data`: adrb-aircraft, feed - wszystkie oznaczone jako `synthetic: true`

### 2.5 Panel administracyjny
- Na podstawie dostępnej dokumentacji: korzysta z HTTP polling i WSS
- Mechanizmy synchronizacji przychodzące z centrali

### 2.6 Mechanizmy cache, deduplikacji i agregacji
- Na podstawie dostępnej dokumentacji: podstawowa struktura danych z polami identyfikującymi
- Brak zaawansowanej logiki deduplikacji w publikowanych fragmentach
- Formaty danych pozwalają na identyfikację po `id` + `source` + `published_at`

## 3. Rozwiązania nadające się do bezpiecznego wykorzystania jako oddzielny moduł

### 3.1 Możliwe do wydobycia bez naruszania prywatności
- Schematy formatu aircraft/airspace/official_notice
- Ogólny pomysł adapterów źródłowych
- Koncepcja lokalnego panelu WWW
- Formaty API (REST + WebSocket)

### 3.2 Unikaj kopiowania
- Kontrakty protokołów WSS (uwierzytelnianie, specyficzne komunikaty)
- Dane użytkowników i konfiguracje VPS
- Klucze API i tokeny dostępowe
- Architekturę bazodanową specyficzną dla produkcji

## 4. Weryfikacja licencji

- **Wersje do 1.7.53**: MIT - można wykorzystać koncepcje, ale nie kopiować kodu
- **Od 1.7.54**: Restrykcyjna licencja - absolutnie zakazane kopiowanie kodu
- **Decyzja**: Pełna independently implementation w Go od zera
- **Licencja projektu**: MIT dla nowego oryginalnego kodu
- **Zależności zewnętrzne**: Weryfikacja każdej biblioteki pod kątem licencji

## 5. Formaty danych i normalizacja (z istniejących fragmentów)

Wspólny mianownik obserwowanych typów:

| Typ | Identyfikator | Źródło | Czas | Pozycja | Wysokość | Parametry | Pewność | Licencja |
|-----|--------------|--------|------|---------|----------|-----------|---------|----------|
| aircraft | hex + source | source | seen | lat, lon | alt_baro | gs, track | manual | source-specific |
| helicopter | hex + source | source | seen | lat, lon | alt_baro | gs, track | manual | source-specific |
| balloon | hex + source | source | seen | lat, lon | alt_baro | gs, track | manual | source-specific |
| satellite | NORAD/CAT | source | published | lat, lon | altitude | - | computed | orbital data license |
| airspace | polygon/rect | source | valid_from/to | - | - | - | computed | licensing-dependent |
| official_notice | id + source | source | published_at | - | - | category, region | manual | authority-specific |
| media_report | id + source | source | published_at | optional lat/lon | - | source_url | manual | media license |
| community_observation | id + source | user | observed_at | lat, lon | - | - | user-provided | community-specific |
| sensor_detection | sensor_id | sensor | detected_at | lat, lon | - | sensor params | computed | sensor-specific |

## 6. Wnioski z audytu

1. **Niezależność**: NieboGuard Node musi być całkowicie samodzielny - nie może polegać na kodzie prywatnego repozytorium
2. **Kompatybilność**: Musi rozumieć te same formaty danych (aircraft, helicopters, etc.) by wymieniać się danymi z centralą
3. **Protokół**: API REST + WebSocket, identyczne nazewnictwo endpointów
4. **Adaptery**: Modułowy system źródeł danych - konieczny oddzielny moduł dla każdego typu źródła
5. **Licencja**: Tylko kod do wersji 1.7.53 (MIT) może być punktem odniesienia koncepcyjnym, nie kopiować bezpośrednio
6. **Baza danych**: SQLite z WAL wydaje się naturalnym wyborem biorąc pod uwagę lokalny charakter projektu
7. **Panel WWW**: Musi działać całkowicie lokalnie, bez wymogu połączenia z centralą
8. **Tryby pracy**: LOCAL (pełna funkcjonalność offline) jako podstawa, COMMUNITY i PUBLIC NODE jako rozszerzenia opcjonalne

## 7. Rekomendacje

1. Rozpocząć od tworzenia dokumentacji architektonicznej w tym repozytorium
2. Zaimplementować pilotowy moduł adaptera ADS-B w Go
3. Zbudować podstawowy silnik REST API + WebSocket
4. Stworzyć panel WWW w HTML/TS obsługujący tryb LOCAL
5. Zaprojektować system normalizacji danych zgodnie ze schematami zidentyfikowanymi w audycie
6. Przygotować instalatory Windows/Linux i konfigurację Podman Quadlet równolegle
7. Unikać jakichkolwiek kodów źródłowych z wersji NieboGuard po 1.7.54