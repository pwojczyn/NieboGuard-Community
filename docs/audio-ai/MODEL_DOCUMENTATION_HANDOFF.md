# Audio AI — handoff dokumentacji modelu

Data: 2026-11-04
Adresat: agent/LLM pracujący w `pwojczyn/NieboGuard-Community`
Status: **materiał do przygotowania dokumentacji; własny model jeszcze nie istnieje.**

Ten plik nie jest model card ani zapowiedzią publikacji. Opisuje fakty,
wymagane dane i szablon dokumentacji, którą Community ma przygotować, gdy
właściciel przekaże wytrenowany model. Nie wymyślaj brakujących wartości.

## 1. Fakty na dziś (2026-11-04)

### Model bazowy — istnieje i jest opublikowany

| Pole | Wartość |
|---|---|
| `id` | `yamnet-base` |
| `version` | `1.0.0` |
| `displayName` | YAMNet base (AudioSet) |
| `baseModel` | `google/yamnet` |
| Format web | ONNX, `model.onnx` |
| Rozmiar ONNX | 16 093 355 B |
| SHA-256 ONNX | `d3835ffbbd4a1bb3e777f0ca217b5007907f5171dd5d17c4236b95b2af8f908e` |
| Format Android | TFLite, `model.tflite` |
| Rozmiar TFLite | 4 126 810 B |
| SHA-256 TFLite | `10c95ea3eb9a7bb4cb8bddf6feb023250381008177ac162ce169694d05c317de` |
| Klasy | 521 klas AudioSet |
| Licencja | Apache-2.0 (TensorFlow Models / AudioSet) |
| Wejście | mono 16 kHz, okno 15 600 próbek, waveform |
| Wyjście | `output_0`, `[frames, 521]`, niezależne scores |
| Kanał | `stable`, manifest podpisany Ed25519 |
| R2 | `audio-ai/models/yamnet-base/1.0.0/` |

YAMNet **nie ma klasy drona wielowirnikowego**. Wyniki typu `Aircraft`,
`Helicopter`, `Propeller` oznaczają podobieństwo akustyczne, nie potwierdzenie
obiektu.

### Własny model — jeszcze nie istnieje

- Nazwa zarezerwowana: `nieboguard-audio-v0.1.0`.
- `id`: `nieboguard-audio`, `version`: `0.1.0`, `displayName`:
  `nieboguard-audio-v0.1.0`, `baseModel`: `google/yamnet`.
- Stan: **`planned`, `approved: false`, brak artefaktu**.
- API `nieboguard-audio/versions` zwraca pustą listę.
- W R2 `audio-ai/models/` nie ma katalogu `nieboguard-audio/`.
- Nie ma wytrenowanych wag, metryk ani datasetu treningowego zatwierdzonego
  do tego modelu.

**Nie pisz, że model istnieje, dopóki właściciel nie przekaże artefaktu.**

## 2. Dane treningowe i licencje — co trzeba udokumentować

Każdy model musi mieć tabelę danych: identyfikator, wersję, klasy, liczbę
nagrań, licencję, źródło, prawo do treningu, prawo do redystrybucji i dowód
licencji. Brak `license_evidence` = brak treningu.

### Znane zbiory

| Dataset | Wersja | Klasy | Licencja | Źródło | Trening | Redystrybucja | Stan |
|---|---|---|---|---|---|---|---|
| `svanstrom-audio` | v1 | background, drone_multirotor, helicopter | **CC0-1.0** | GitHub `DroneDetectionThesis/Drone-detection-dataset` | tak | tak | 100 shardów (wg docs) |
| `svanstrom-audio` | v2 | j.w. | **CC0-1.0** | j.w. | tak | tak | 84 shardy (wg docs) |
| `admin` (nagrania terenowe) | — | background / drone / helicopter | **DO POTWIERDZENIA** | pliki administracyjne w R2 | **nie zakładać** | **nie zakładać** | 394 WAV, ok. 55 MB; pochodzenie i zgody wymagają weryfikacji |
| `uavirbase` | DRAFT | — | CC-BY-4.0 (potwierdzona w audycie) | — | tak | tak | brak manifestu/SHA; nieużywany |
| `fsd50k` | — | — | wymaga weryfikacji per klip | Freesound | nie | nie | `BLOCKED` |
| dane syntetyczne testowe | — | — | CC0 / synthetic | skrypty testowe | tylko testy | nie dotyczy | nie są zbiorem produkcyjnym |

Zasady:

- Split **po `source_recording_id`**, nie po oknie — wszystkie okna jednego
  nagrania w jednym zbiorze.
- Test musi pochodzić z innych nagrań/urządzeń niż trening.
- Nie używaj `admin/` w modelu publicznym, dopóki nie ma udokumentowanej
  licencji/zgody i dowodu.
- Nagrania wolontariuszy wymagają świadomej zgody i informacji o wykorzystaniu.
- Zachowaj atrybucję wymaganą przez licencję i link do tekstu licencji.
- Nie publikuj surowych nagrań ani danych osobowych w repo.

### Szablon tabeli do wypełnienia

```markdown
| Dataset | Version | Classes | Recordings | Duration | License | License URL | Training allowed | Redistribution | Evidence |
|---|---|---|---:|---:|---|---|---|---|---|
| <id> | <v> | <klasy> | <n> | <h> | <licencja> | <url> | yes/no | yes/no | <plik/URL> |
```

## 3. Licencja własnego modelu — decyzja właściciela

Licencja kodu Community (MIT) **nie jest automatycznie licencją wag modelu**.
Właściciel musi wybrać i zapisać jedną z opcji:

| Opcja | Komercyjne użycie | Redystrybucja wag | Wymagania |
|---|---|---|---|
| CC BY 4.0 | tak | tak | atrybucja, link do licencji |
| CC0 1.0 | tak | tak | brak obowiązkowej atrybucji (dobra wola) |
| Apache-2.0 | tak | tak | NOTICE, klauzula patentowa |
| MIT | tak | tak | nota copyright |
| CC BY-NC 4.0 | nie | tak | tylko niekomercyjne; sprawdzić zgodność z celem projektu |
| Własna licencja | zależy | zależy | tekst musi być w repo; nie może być sprzeczny z licencjami danych |
| Model zamknięty | zależy | **nie** | nie publikować wag; tylko opis i ewentualnie API |

Wymagane ustalenia:

- `<MODEL_LICENSE>` — wybrana licencja wag,
- `<MODEL_LICENSE_URL>` — link lub plik w repo,
- `<MODEL_COPYRIGHT>` — właściciel i rok,
- `<BASE_MODEL_NOTICE>` — Apache-2.0 dla `google/yamnet`,
- `<DATASET_ATTRIBUTION>` — atrybucje wymagane przez zbiory.

Jeśli właściciel nie podejmie decyzji, wpisz `NEEDS CONFIRMATION` i nie
publikuj wag ani linku do artefaktu.

## 4. Model card — wymagane sekcje

Po otrzymaniu modelu utwórz `docs/audio-ai/MODEL_CARD.md`:

1. **Identyfikacja** — `id`, `version`, `displayName`, `baseModel`,
   `parentVersion`, `status`, `approved`, data, autor.
2. **Przeznaczenie** — do czego model służy i czego nie robi; wyraźnie:
   brak automatycznego alarmu, brak potwierdzenia zagrożenia.
3. **Architektura i wejście/wyjście** — typ, rozmiar, sample rate, klasy,
   kolejność klas, preprocessing, progi.
4. **Dane treningowe** — tabela z sekcji 2, split, liczność, źródła,
   pochodzenie, zgodność licencji, data pobrania.
5. **Trening** — commit kodu, framework, hiperparametry, seed, sprzęt, czas,
   wersja datasetu, reprodukcja.
6. **Ewaluacja** — metryki per klasa, macierz pomyłek, FAPH/FA/h, zbiór testowy
   niezależny od treningu; brak metryk = brak publikacji.
7. **Ograniczenia** — dźwięki tła, wiatr, ruch drogowy, urządzenia, język,
   klasy szerokie; YAMNet nie ma klasy drona.
8. **Prywatność i etyka** — brak surowych nagrań w repo, zgody, dane osobowe,
   przeznaczenie cywilne, brak nadzoru.
9. **Dystrybucja** — R2 key, rozmiar, SHA-256, podpis, format, platformy,
   kanał, wymagana wersja aplikacji.
10. **Licencja i atrybucja** — sekcja 3 + licencje danych + NOTICE.
11. **Historia wersji** — `parentVersion`, data, zmiany, status.

## 5. Manifest — szkic do uzupełnienia

```json
{
  "schema": "nieboguard.audio-model/1",
  "id": "nieboguard-audio",
  "version": "0.1.0",
  "displayName": "nieboguard-audio-v0.1.0",
  "baseModel": "google/yamnet",
  "parentVersion": null,
  "status": "candidate",
  "approved": false,
  "format": "onnx",
  "platforms": ["web", "android"],
  "sampleRate": 16000,
  "windowSamples": 15600,
  "classes": ["<kolejność klas>"],
  "preprocess": { "type": "waveform", "normalize": "none", "shape": [15600] },
  "postprocess": { "activation": "none" },
  "artifact": {
    "url": "<R2/Worker URL>",
    "size": 0,
    "sha256": "<64 hex>"
  },
  "artifacts": {
    "android": { "url": "<R2/Worker URL>", "size": 0, "sha256": "<64 hex>" }
  },
  "license": "<MODEL_LICENSE>",
  "datasets": ["<dataset-id>@<version>"],
  "metrics": { "per_class": {}, "faph": null },
  "publishedAt": null,
  "signature": "<Ed25519>"
}
```

Nie wstawiaj `approved: true` ani `publishedAt` bez decyzji człowieka.

## 6. Publikacja z opóźnieniem — workflow

1. Właściciel przekazuje artefakt, manifest, metryki i listę datasetów
   (prywatny kanał).
2. Weryfikacja: SHA-256, podpis, klasy, preprocessing, metryki, licencje.
3. Niezależna akceptacja człowieka; dopiero wtedy `approved: true`.
4. Community publikuje **najpierw dokumentację** (model card + licencje),
   bez linku do artefaktu, jeśli właściciel jeszcze go nie udostępnił.
5. Po decyzji o publikacji: manifest w R2, kanał `stable`/`beta`, link w
   README, changelog i release notes.
6. Nie publikuj automatycznie do Google Play ani GitHub Releases.
7. Nie ogłaszaj modelu dronów bez metryk i testów terenowych.

## 7. Pliki do utworzenia w Community

```text
docs/audio-ai/README.md
docs/audio-ai/MODEL_CARD.md
docs/audio-ai/DATASETS_AND_LICENSES.md
docs/audio-ai/MODEL_LICENSE.md
docs/audio-ai/PUBLICATION_CHECKLIST.md
```

README główny Community: dodać sekcję „Audio AI” i link do powyższych.
Nie kopiuj do Community prywatnych kluczy, tokenów R2, adresów VPS ani
surowych danych.

## 8. Zadanie dla agenta Community (kopiuj/wklej)

> Przeczytaj `docs/audio-ai/MODEL_DOCUMENTATION_HANDOFF.md`. Nie wymyślaj
> danych. Utwórz szkielet dokumentacji z sekcji 7, oznacz każdą brakującą
> wartość jako `NEEDS CONFIRMATION`, nie publikuj linku do modelu i nie
> twierdź, że `nieboguard-audio-v0.1.0` istnieje. Użyj faktów o `yamnet-base`
> z sekcji 1. Przygotuj PR z dokumentacją i listą pytań do właściciela o
> licencję wag, dataset `admin/` i metryki. Uruchom testy linków/markdown i
> `python tests/test_repo_hygiene.py`, jeśli istnieje.

## 9. Checklist akceptacji dokumentacji

- [ ] Nie ma twierdzenia, że własny model istnieje przed publikacją.
- [ ] `yamnet-base` jest opisany jako model bazowy, nie detektor dronów.
- [ ] Jest tabela datasetów z licencjami i dowodami.
- [ ] `admin/` ma status „do potwierdzenia” i nie jest domyślnie treningowy.
- [ ] Jest wybrana licencja wag albo `NEEDS CONFIRMATION`.
- [ ] Jest NOTICE dla Apache-2.0 YAMNet.
- [ ] Są metryki per klasa i informacja o zbiorze testowym.
- [ ] Brak surowych nagrań, danych osobowych, kluczy i tokenów.
- [ ] Linki działają; brak obietnic Play/Community release.
