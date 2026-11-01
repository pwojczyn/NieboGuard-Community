# Sensors — RSSI i odległość

`tools/sensors/rssi_distance.py` szacuje odległość z RSSI i modelu log-distance.
To **bardzo przybliżone** i zależy od anteny, mocy, przeszkód i orientacji.

```bash
python tools/sensors/rssi_distance.py --rssi -70 --tx-power -59 --n 2.0
```

## Ograniczenia

- nie używaj do precyzyjnej lokalizacji,
- nie publikuj współrzędnych prywatnych odbiorników,
- nie przedstawiaj wyniku jako potwierdzenia obecności drona,
- kalibracja wymaga własnych pomiarów w kontrolowanych warunkach.

## Test

```bash
python -m unittest tests.test_rssi_distance -v
```
