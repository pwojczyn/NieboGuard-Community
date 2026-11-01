# Audio AI tools

## `wav_audit.py`

Reads WAV header metadata without uploading or decoding the signal:

```bash
python tools/audio-ai/wav_audit.py sample.wav
```

Returns duration, sample rate, channels, bit depth and warnings. No network,
no model, no automatic upload.

## Label schema

[`labels.example.yaml`](labels.example.yaml) is a starting vocabulary. Multiple
labels are allowed. `unknown` and `NOT_SURE` are valid. A model suggestion is
not a training label until a human accepts it.

Do not publish private conversations or samples without consent and a clear
license. YAMNet is a broad AudioSet classifier, not a drone detector.

## NieboGuard Compute

Trening małej głowy klasyfikatora na cechach audio jest w
[`tools/nieboguard-compute`](../nieboguard-compute/README.md). Klient nie
pobiera całych archiwów; używa małych paczek z manifestem i SHA-256.
