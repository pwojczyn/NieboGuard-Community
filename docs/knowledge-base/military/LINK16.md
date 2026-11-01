# Link 16 (STANAG 5516)

## Introduction (Polski)
Link 16 to międzynarodowy standard STANAG 5516, który definiuje system łączności danych dla wiarygodnych systemów obrony powietrznej i przeciwlotniczej. Jest to fundamentowy protokół dla systemów A2/AD (Air-to-Air/air-to-ground) i systemów radarowych.

## Introduction (English)
Link 16 (STANAG 5516) is an international standard that defines a data link system for air defense and surveillance systems. It is the fundamental protocol for air-to-air/air-to-ground systems and radar systems.

## Architecture (Polski)
- **Warstwa fizyczna**: VHF/UHF (225-400 MHz) z modulacją FM lub PM
- **Warstwa danych**: Pakiety STANAG 5516/STANAG 4601
- **Warstwa aplikacji**: Systemy A2/AD, systemy radarowe, systemy EW

## Architecture (English)
- **Physical layer**: VHF/UHF (225-400 MHz) with FM or PM modulation
- **Data layer**: STANAG 5516/STANAG 4601 packets
- **Application layer**: A2/AD systems, radar systems, EW systems

## Data Types (Polski)
- **Mode 1**: Statusy systemów, informacje o stanie
- **Mode 2**: Pozycje, wektory prędkości, informacje o lokalizacji
- **Mode 3**: Współrzędne, informacje o celach, dane o lokalizacji
- **Mode 4**: Informacje o systemach, statusy, dane o systemach

## Data Types (English)
- **Mode 1**: System statuses, status information
- **Mode 2**: Positions, velocity vectors, location information
- **Mode 3**: Coordinates, target information, location data
- **Mode 4**: System information, statuses, system data

## Communication Protocol
- **Format**: Pakiety STANAG 5516/STANAG 4601
- **Struktura pakietu**:
  - Nagłówek (16 bitów)
  - Tryb (4 bity)
  - Długość danych (12 bitów)
  - Dane (długość określona w nagłówku)
  - CRC (16 bitów)
- **Prędkość**: 50-100 kbps (zależnie od trybu)
- **Zasięg**: 200-400 km (zależnie od warunków)

## Limitations (Polski)
- Ograniczone pole widzenia (line-of-sight)
- Zależność od warunków atmosferycznych
- Ograniczona liczba węzłów w sieci
- Brak zabezpieczeń kryptograficznych w wersji podstawowej

## Limitations (English)
- Line-of-sight limited range
- Atmospheric condition dependent
- Limited node count in network
- No cryptographic security in basic version

## Implementation Status
- **NieboGuard**: Adapter stub (planned)
- **Simulator**: Educational synthetic data generator (planned)
- **Status**: PLANNED - requires adapter development

## Integration with NieboGuard
- Adapter will connect to Link 16 data sources (simulated or real)
- Data will be processed through NieboGuard's core engine
- WebSocket updates will be available for panels
- REST API will provide access to Link 16 data

## Resources
- STANAG 5516 (official)
- STANAG 4601 (supplemental)
- NATO STANAG documentation
- Public domain technical specifications

## Status
PLANNED - Adapter and simulator development required