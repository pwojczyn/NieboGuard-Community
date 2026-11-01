# Architektura — widok publiczny

```text
publiczne źródła danych
        │
        ▼
NieboGuard (prywatny backend)
        │  publiczne endpointy read-only
        ▼
SDK / narzędzia Community / Twoja aplikacja
```

## Komponenty

- **Frontend publiczny** — mapa, komunikaty, alarmy, źródła, statystyki i wiedza.
- **Publiczne API read-only** — stan, feed, alarmy, incydenty, źródła,
  statystyki i lokalizacje.
- **SDK Community** — cienka warstwa HTTP bez sekretów i bez zapisu.
- **Narzędzia** — ADS-B → GeoJSON, audyt WAV, audyt obrazów, szacowanie RSSI.
- **Dane przykładowe** — wyłącznie syntetyczne, do testów.

## Granice publiczne

- Nie publikujemy adresów wewnętrznych VPS, konfiguracji WSS, schematu D1,
  tokenów ani kodu administracyjnego.
- Publiczne API nie wymaga tokenu do odczytu; endpointy zapisu i administracji
  są poza zakresem Community.
- Dane mogą być opóźnione, niepełne lub niepotwierdzone.
- NieboGuard **nie jest** oficjalnym systemem ostrzegania.

## Zasady projektowe

1. Najpierw źródło i licencja, potem dane.
2. Brak danych nie oznacza braku zagrożenia.
3. Alarm/ostrzeżenie nie jest automatycznie incydentem.
4. Dane syntetyczne są jawnie oznaczane.
5. Wolontariusze nie publikują prywatnych lokalizacji ani surowych nagrań.
