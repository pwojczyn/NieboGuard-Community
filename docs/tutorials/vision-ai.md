# Vision AI — bezpieczny start

## Zasady

- używaj obrazów, do których masz prawa,
- usuń EXIF/GPS przed publikacją,
- nie publikuj twarzy, tablic ani prywatnych instalacji,
- nie wymuszaj identyfikacji modelu na niewyraźnym zdjęciu,
- `NOT_SURE` jest poprawną etykietą.

## Audyt obrazu

```bash
python tools/vision-ai/image_audit.py photo.jpg
python tools/vision-ai/image_audit.py photo.jpg --strip /tmp/photo-clean.jpg
```

Pierwsze polecenie pokazuje wymiary, tryb i obecność EXIF. Drugie zapisuje
kopię bez metadanych (wymaga Pillow).

## Adnotacje

- prostokąty i klasy powinny być opisane w manifeście,
- nie mieszaj tego samego nagrania w train i test,
- brak detekcji nie jest potwierdzeniem braku obiektu,
- wynik modelu trafia do review, nie do publicznego alarmu.

## Licencje

Sprawdź licencję obrazów i adnotacji. Open Images/FSD50K/Roboflow mają różne
warunki per element; nie zakładaj, że „open” oznacza dowolne użycie.
