# Audio AI — bezpieczny start

## Zasady

- nagrywaj tylko dźwięki, do których masz prawo,
- nie publikuj prywatnych rozmów ani danych osobowych,
- nie wysyłaj automatycznie mikrofonu,
- każda próbka wymaga świadomej zgody i informacji o wykorzystaniu,
- wynik modelu nie jest potwierdzeniem zagrożenia.

## Audyt WAV

```bash
python tools/audio-ai/wav_audit.py sample.wav
```

Narzędzie czyta nagłówek WAV i zwraca czas, sample rate, kanały i szerokość
próbki. Nie wysyła pliku.

## Etykiety

Przykładowy schemat: [`tools/audio-ai/labels.example.yaml`](../../tools/audio-ai/labels.example.yaml).

- wiele etykiet na nagranie jest dozwolone,
- `unknown` i `NOT_SURE` są poprawne,
- automatyczna sugestia modelu nie jest etykietą treningową,
- YAMNet nie ma klasy drona; nie przedstawiaj go jako detektora dronów.

## Licencje

Nie dodawaj datasetów bez sprawdzenia licencji audio, adnotacji i warunków API.
Freesound nie jest automatycznie komercyjnie dozwolony.
