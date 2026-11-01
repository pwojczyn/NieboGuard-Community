# Sensor tools

`rssi_distance.py` is a rough log-distance estimate:

```bash
python tools/sensors/rssi_distance.py --rssi -70 --tx-power -59 --n 2.0
```

This is **not** precise localization and not proof that a drone is present.
Calibrate for your own antenna/environment. Never publish private receiver
coordinates.
