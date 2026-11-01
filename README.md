# NieboGuard Community

Otwarte dokumenty, SDK i narzędzia dla społeczności wokół [NieboGuard](https://nieboguard.pl/).

> **NieboGuard nie jest oficjalnym systemem ostrzegania ludności.** Nie zastępuje
> Alertów RCB, syren, komunikatów służb ani numeru 112. Dane mogą być opóźnione,
> niepełne lub niepotwierdzone.

## Co jest w tym repozytorium

- dokumentacja publicznego API i architektury,
- SDK JavaScript i Python do publicznych endpointów read-only,
- narzędzia ADS-B, audio, vision i sensorów,
- **NieboGuard Compute** — klient rozproszonego trenowania Audio AI/Vision AI,
  instalatory Windows/Linux, małe paczki danych i audyt źródeł akademickich,
- syntetyczne dane testowe,
- tutoriale dla wolontariuszy,
- szablony zgłoszeń błędów i propozycji funkcji.

## Czego tu nie ma

- kodu produkcyjnego NieboGuard,
- sekretów, konfiguracji VPS, kluczy API, APK,
- danych użytkowników i surowych lokalizacji,
- prywatnej historii Git ani wewnętrznych dokumentów operacyjnych.

Główne repozytorium produkcyjne pozostaje prywatne. To repozytorium jest
publicznym zapleczem dla dokumentacji, narzędzi i współpracy.

## Szybki start

```bash
git clone https://github.com/pwojczyn/NieboGuard-Community.git
cd NieboGuard-Community
python3 -m venv .venv
. .venv/bin/activate
python -m unittest discover -s tests -v
node --test tests/*.test.mjs
```

SDK JavaScript nie ma zależności produkcyjnych. SDK Python i klient Compute
korzystają wyłącznie ze standardowej biblioteki. Narzędzie vision opcjonalnie
wymaga Pillow.

## NieboGuard Compute

Klient społecznościowego trenowania AI jest w
[`tools/nieboguard-compute/`](tools/nieboguard-compute/README.md). Obsługuje
Windows/Linux, CPU i opcjonalne NVIDIA GPU, małe paczki danych MINI/STANDARD/
EXTENDED/GPU oraz bezpieczny safetensors. Udział wymaga jawnej zgody; klient
nie działa w tle i nie wykonuje kodu z backendu.

```bash
cd tools/nieboguard-compute
python -m unittest discover -s tests -v
```

## Testowy APK

W [`releases/`](releases/README.md) jest wyłącznie testowy build NieboGuard:

- `nieboguard-mobile-0.1.0-dev-debug.apk`,
- podpis debug, nie release,
- SHA-256 w [`releases/SHA256SUMS`](releases/SHA256SUMS).

To jedyny plik APK w tym repozytorium. Nie ma tu `Straznik.apk` ani innych
binarek. APK nie zawiera sekretów serwerowych; szczegóły w
[`releases/APK-NOTICE.md`](releases/APK-NOTICE.md).

## Publiczne API

- dokumentacja: [`docs/api/README.md`](docs/api/README.md)
- OpenAPI: [`docs/api/openapi.yaml`](docs/api/openapi.yaml)
- przykład: [`examples/state_watch.mjs`](examples/state_watch.mjs)

Publiczne endpointy są tylko do odczytu i nie wymagają tokenu. Nie wysyłaj do
nich danych osobowych, współrzędnych prywatnych instalacji ani sekretów.

## Jak pomóc

- popraw dokumentację lub tłumaczenie,
- dodaj test lub narzędzie,
- zgłoś błąd przez GitHub Issues,
- opisz przypadek użycia odbiornika ADS-B, audio lub vision,
- zgłoś lukę bezpieczeństwa prywatnie według [`SECURITY.md`](SECURITY.md).

Zasady: [`CONTRIBUTING.md`](CONTRIBUTING.md) · [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).

## Licencja

Nowe pliki tego repozytorium są udostępniane na licencji MIT
([`LICENSE`](LICENSE)). Pochodzenie kodu bazowego Strażnik i zależności opisuje
[`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md).

To repozytorium **nie** udziela licencji do produkcyjnego kodu NieboGuard,
danych użytkowników, APK ani infrastruktury.
