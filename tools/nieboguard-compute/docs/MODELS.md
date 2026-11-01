# Audyt modeli Audio AI i Vision AI

Data: 2026-11-01.

## Co faktycznie działa

| Moduł | Model | Klasy | Gdzie | Status |
|---|---|---|---|---|
| Android Audio Scanner | YAMNet TFLite | 521 klas AudioSet; brak klasy drona | urządzenie | działa lokalnie |
| `acoustic-training/train_yamnet.py` | MLP na embeddingach YAMNet | konfigurowalne | Python/TensorFlow | kod istnieje, brak wytrenowanego modelu |
| Compute Audio AI | mała głowa na cechach log-mel lub embeddingach | dron, samolot, śmigłowiec, ptak, tło | klient CPU/GPU | wytrenowana na syntetycznym fixturze; brak realnego modelu |
| `analyzeImage` | Workers AI / OpenRouter vision LLM | opis obrazu | Worker | nie jest detektorem, nie zwraca bbox |
| Vision Lab | brak modelu | etykiety ręczne | panel | brak detekcji |
| Compute Vision AI | głowa detekcji na cechach propozycji | dron, samolot, śmigłowiec, ptak, balon | klient CPU/GPU | syntetyczny test; brak realnego modelu |
| Sky Lab | symulator | scenariusze | frontend | nie jest modelem AI |

## Czego nie ma

- wytrenowanego detektora YOLO dla małych obiektów na tle nieba,
- realnego modelu dronów na danych akademickich,
- niezależnego holdoutu z prawdziwych danych,
- serwerowej inferencji audio/wideo,
- automatycznej detekcji w Vision Lab,
- potwierdzonej jakości klas dron/śmigłowiec/ptak/balon.

## YAMNet — ograniczenia

YAMNet klasyfikuje szerokie kategorie AudioSet. Nie ma klasy „dron”. Może
wykryć „propeller”, „aircraft”, „vehicle”, „wind”, ale nie jest detektorem
dronów. Trening głowy na embeddingach YAMNet jest sensownym startem, ale
wymaga danych z dronami i niezależnej walidacji.

## Vision AI — plan

1. Lekki backbone (np. MobileNet/SSDlite lub mały YOLO) zamrożony.
2. Głowa detekcji: objectness + klasy + box.
3. Klasy: dron, samolot, śmigłowiec, ptak, balon.
4. Małe obiekty: augmentacja, kafelkowanie, próg pewności, analiza FP.
5. Eksport safetensors/ONNX/TFLite po walidacji.
6. Nie deklarować identyfikacji typu ani zagrożenia z obrazu.

Pełny trening backbone od zera jest poza MVP. Compute trenuje głowę na
cechach/propozycjach; opcjonalny PyTorch może użyć GPU.

## Zasady publikacji

- Brak metryk = brak dowodu jakości.
- Raport musi zawierać accuracy, base_accuracy, false alarm rate i per-class
  recall.
- Rollback jest obowiązkowy.
- Nie publikujemy modelu, jeśli dataset nie ma prawa redystrybucji i treningu.
- Nie twierdzimy, że model rozpoznaje drony na podstawie samego YAMNet.
