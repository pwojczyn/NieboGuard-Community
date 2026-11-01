# Akademickie źródła danych — audyt

Data: 2026-11-01. Sprawdzono publiczne API/repozytoria. To rejestr
operacyjny, nie opinia prawna. Żaden z tych zbiorów nie jest jeszcze
udostępniany wolontariuszom; wpisy `DRAFT` wymagają manifestu, SHA-256,
filtrowania i przeglądu praw.

## Podsumowanie

| Źródło | Autorzy | Licencja | Rozmiar | Status | Redystrybucja |
|---|---|---|---|---|---|
| Svanström multi-sensor | F. Svanström, F. Alonso-Fernandez, C. Englund | CC0-1.0 | 90 audio, 650 wideo, 203 328 obrazów | `DRAFT` | tak (CC0) |
| Politechnika Świętokrzyska | M. Mięsikowska | CC-BY-4.0 | ~12 GB, 17 archiwów | `DRAFT` | tak, po odfiltrowaniu mowy |
| WAT UaVirBASE | G. Jekateryńczuk, R. Szadkowski, Z. Piotrowski | CC-BY-4.0 | 13,8 GB | `DRAFT` | tak, z atrybucją |
| University of Surrey DDL | S. Safavi, M. A. Safavi, B. Cook, W. Wang | CC-BY-4.0 | 17,27 GB | `DRAFT` | tak, z atrybucją |
| YOLOBirDrone | autorzy / CSIR | dostęp na wniosek; brak jasnej licencji redystrybucji | nieznany | `BLOCKED` | **nie** bez pisemnej zgody |

## 1. Drone Detection Dataset — Fredrik Svanström

- Repozytorium: <https://github.com/DroneDetectionThesis/Drone-detection-dataset>
- Licencja: **CC0-1.0** (plik LICENSE w repozytorium).
- Zawartość: 90 klipów audio, 650 wideo (365 IR, 285 visible), 203 328
  opisanych obrazów; klasy wideo: Airplane, Bird, Drone, Helicopter; audio:
  Drone, Helicopter, Background.
- Ograniczenia: adnotacje `.mat` wymagają dekodera `mcos-decoder`; duże
  archiwum wideo; trzeba wybrać małe paczki i podzielić po sesjach/klipach.
- Decyzja: **kandydat DRAFT**. Nie publikować całych archiwów. CC0 pozwala na
  redystrybucję, ale i tak wymagamy atrybucji w metadanych.

## 2. Politechnika Świętokrzyska — nagrania akustyczne UAV

- Zenodo: <https://zenodo.org/records/15190811>
- DOI: `10.5281/zenodo.15190811`
- Autor: Marzena Mięsikowska (Kielce University of Technology).
- Licencja: **CC-BY-4.0** (API Zenodo).
- Zawartość: 17 archiwów X4/X6, około 12 GB; nagrania dronów i mowy w
  warunkach zewnętrznych.
- Ograniczenia: **mowa ludzka** — nie może trafić do paczek wolontariuszy.
  Wymagane filtrowanie, weryfikacja braku danych osobowych i manifest.
- Decyzja: **kandydat DRAFT**, tylko po filtrze mowy.

## 3. Wojskowa Akademia Techniczna — UaVirBASE

- Zenodo: <https://zenodo.org/records/15391924>
- DOI: `10.5281/zenodo.15391924`
- Autorzy: Gabriel Jekateryńczuk, Rafał Szadkowski, Zbigniew Piotrowski.
- Licencja: **CC-BY-4.0**.
- Zawartość: 13,8 GB, `Microphone_array.zip`; synchronizowane nagrania
  wielomikrofonowe do lokalizacji źródła dźwięku UAV (SSL).
- Ograniczenia: to głównie SSL, nie klasyfikacja dron/nie-dron; potrzebne
  mapowanie etykiet i sprawdzenie, czy da się zbudować klasyfikator.
- Decyzja: **kandydat DRAFT**. Duże archiwum importuje zewnętrzny runner,
  nie Worker.

## 4. University of Surrey — DDL

- Zenodo: <https://zenodo.org/records/6459183>
- DOI: `10.5281/zenodo.6459183`
- Autorzy: Saeid Safavi, Mohammad Amin Safavi, Ben Cook, Wenwu Wang.
- Licencja: **CC-BY-4.0**.
- Zawartość: 17,27 GB, `MLSP_2022_Real_Data.zip` i
  `MLSP_2022_Synth_Data.zip`; dane audio/mikrofonowe do detekcji i lokalizacji.
- Ograniczenia: duże archiwa, Real/Synth; trzeba potwierdzić brak mowy i
  podzielić po sesjach; osobny manifest dla real/synth.
- Decyzja: **kandydat DRAFT**.

## 5. YOLOBirDrone

- arXiv: <https://arxiv.org/abs/2601.08319>
- Repozytorium: <https://github.com/dapinderk-2408/YOLOBirDrone>
- Licencja: **brak jednoznacznej licencji**; „Available on Request”, formularz
  i kontakt do CSIR.
- Decyzja: **BLOCKED** dla redystrybucji wolontariuszom. Można użyć wyłącznie
  po pisemnej zgodzie i przeglądzie prawnym. Nie publikujemy go w paczkach.

## Zasady paczek

- Tylko źródła z `training_allowed=1` i `redistribution_allowed=1`.
- Każdy shard ma SHA-256, grupę sesji i split.
- Model i biblioteki mają osobny limit od paczki danych.
- Nie dzielimy jednej sesji między train/validation/test ani między paczki
  tej samej klasy.
- Nie publikujemy surowych archiwów; budujemy małe paczki z manifestem.
- Atrybucja i URL źródła są w manifescie zadania.

## Czego jeszcze brakuje

- pobrania i przeliczenia SHA-256 rzeczywistych plików,
- manifestów paczek dla tych źródeł,
- filtrowania mowy i danych osobowych,
- decyzji prawnej o YOLOBirDrone,
- niezależnego holdoutu dla każdego źródła.
