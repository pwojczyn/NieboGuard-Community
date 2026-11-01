"""Audio AI: WAV decoding and a small log-mel feature extractor.

The production path can use YAMNet embeddings supplied by the task or computed
locally with TensorFlow (optional). This standard-library path guarantees that
a CPU-only volunteer can fine-tune a small classifier without downloading an
unverified runtime.
"""
from __future__ import annotations

import array
import json
import math
import struct
import wave
from pathlib import Path
from typing import Any

SAMPLE_RATE = 16000
FRAME = 400
HOP = 160
FFT_SIZE = 512
MEL_BINS = 32
FEATURE_DIM = MEL_BINS * 2
MAX_RECORDS = 200000


class AudioError(ValueError):
    pass


def _read_wav(path: Path) -> list[float]:
    try:
        with wave.open(str(path), "rb") as wav:
            channels = wav.getnchannels()
            width = wav.getsampwidth()
            rate = wav.getframerate()
            frames = wav.readframes(wav.getnframes())
    except Exception as exc:
        raise AudioError(f"Nie można odczytać WAV: {exc}") from exc
    if width == 1:
        samples = [(b - 128) / 128.0 for b in frames]
    elif width == 2:
        raw = array.array("h")
        raw.frombytes(frames)
        samples = [v / 32768.0 for v in raw]
    elif width == 3:
        samples = []
        for i in range(0, len(frames) - 2, 3):
            v = frames[i] | (frames[i + 1] << 8) | (frames[i + 2] << 16)
            if v & 0x800000:
                v -= 1 << 24
            samples.append(v / 8388608.0)
    elif width == 4:
        raw = array.array("i")
        raw.frombytes(frames)
        samples = [v / 2147483648.0 for v in raw]
    else:
        raise AudioError("Nieobsługiwana szerokość próbki WAV.")
    if channels > 1:
        mono = []
        for i in range(0, len(samples) - channels + 1, channels):
            mono.append(sum(samples[i:i + channels]) / channels)
        samples = mono
    if rate != SAMPLE_RATE and samples:
        ratio = SAMPLE_RATE / float(rate)
        out_len = max(1, int(len(samples) * ratio))
        resampled = []
        for i in range(out_len):
            pos = i / ratio
            left = int(pos)
            frac = pos - left
            right = min(left + 1, len(samples) - 1)
            resampled.append(samples[left] * (1 - frac) + samples[right] * frac)
        samples = resampled
    if not samples:
        raise AudioError("Pusty plik WAV.")
    return samples


def _fft(values: list[float]) -> list[complex]:
    n = len(values)
    if n & (n - 1):
        raise ValueError("FFT size must be a power of two")
    a = [complex(v, 0.0) for v in values]
    j = 0
    for i in range(1, n):
        bit = n >> 1
        while j & bit:
            j ^= bit
            bit >>= 1
        j ^= bit
        if i < j:
            a[i], a[j] = a[j], a[i]
    length = 2
    while length <= n:
        angle = -2 * math.pi / length
        wlen = complex(math.cos(angle), math.sin(angle))
        for i in range(0, n, length):
            w = 1 + 0j
            for k in range(length // 2):
                u = a[i + k]
                v = a[i + k + length // 2] * w
                a[i + k] = u + v
                a[i + k + length // 2] = u - v
                w *= wlen
        length <<= 1
    return a


def _mel(freq: float) -> float:
    return 2595.0 * math.log10(1.0 + freq / 700.0)


def _mel_inv(mel: float) -> float:
    return 700.0 * (10 ** (mel / 2595.0) - 1.0)


def _mel_filters() -> list[list[tuple[int, float]]]:
    low, high = _mel(20.0), _mel(SAMPLE_RATE / 2.0)
    points = [_mel_inv(low + (high - low) * i / (MEL_BINS + 1)) for i in range(MEL_BINS + 2)]
    bins = [int(math.floor((FFT_SIZE + 1) * p / SAMPLE_RATE)) for p in points]
    filters = []
    for m in range(1, MEL_BINS + 1):
        left, center, right = bins[m - 1], bins[m], bins[m + 1]
        pairs = []
        for k in range(max(0, left), min(FFT_SIZE // 2 + 1, right + 1)):
            if k <= center:
                w = (k - left) / max(1, center - left)
            else:
                w = (right - k) / max(1, right - center)
            if w > 0:
                pairs.append((k, w))
        filters.append(pairs)
    return filters


_FILTERS = _mel_filters()


def wav_features(path: Path, input_dim: int = FEATURE_DIM) -> list[float]:
    samples = _read_wav(path)
    if len(samples) < 400:
        samples = samples + [0.0] * (400 - len(samples))
    frames = []
    for start in range(0, max(1, len(samples) - FRAME + 1), HOP):
        chunk = samples[start:start + FRAME]
        if len(chunk) < FRAME:
            chunk += [0.0] * (FRAME - len(chunk))
        windowed = [chunk[i] * (0.54 - 0.46 * math.cos(2 * math.pi * i / (FRAME - 1))) for i in range(FRAME)]
        spectrum = _fft(windowed + [0.0] * (FFT_SIZE - FRAME))
        power = [(abs(spectrum[k]) ** 2) for k in range(FFT_SIZE // 2 + 1)]
        frame = []
        for pairs in _FILTERS:
            energy = sum(power[k] * w for k, w in pairs)
            frame.append(math.log(max(1e-10, energy)))
        frames.append(frame)
    if not frames:
        frames = [[0.0] * MEL_BINS]
    means = [sum(f[i] for f in frames) / len(frames) for i in range(MEL_BINS)]
    stds = [math.sqrt(sum((f[i] - means[i]) ** 2 for f in frames) / len(frames)) for i in range(MEL_BINS)]
    values = means + stds
    if len(values) < input_dim:
        values += [0.0] * (input_dim - len(values))
    return [round(v, 6) for v in values[:input_dim]]


def _safe_child(root: Path, name: str) -> Path:
    base = Path(name).name
    if not base or base in (".", "..") or "/" in name or "\\" in name:
        raise AudioError("Niedozwolona nazwa pliku w manifeście.")
    if Path(base).suffix.lower() in (".py", ".sh", ".ps1", ".bat", ".exe", ".dll", ".so", ".zip", ".tar", ".gz", ".7z"):
        raise AudioError(f"Niedozwolony typ pliku: {base}")
    return root / base


def _normalize_label(value: Any, classes: list[str]) -> str | None:
    if isinstance(value, str) and value in classes:
        return value
    try:
        idx = int(value)
        if 0 <= idx < len(classes):
            return classes[idx]
    except Exception:
        pass
    return None


def load_audio_dataset(data_dir: Path, input_dim: int, classes: list[str]) -> list[dict[str, Any]]:
    files = sorted(p for p in data_dir.iterdir() if p.is_file() and p.suffix.lower() == ".jsonl")
    if not files:
        raise AudioError("Brak pliku JSONL z etykietami audio.")
    rows: list[dict[str, Any]] = []
    for label_file in files:
        for line_no, line in enumerate(label_file.read_text(encoding="utf-8").splitlines(), 1):
            line = line.strip()
            if not line:
                continue
            try:
                item = json.loads(line)
            except Exception as exc:
                raise AudioError(f"{label_file.name}:{line_no}: nieprawidłowy JSON") from exc
            if not isinstance(item, dict):
                raise AudioError(f"{label_file.name}:{line_no}: wiersz nie jest obiektem")
            label = _normalize_label(item.get("label"), classes)
            if label is None:
                raise AudioError(f"{label_file.name}:{line_no}: nieznana etykieta")
            features = item.get("features")
            if features is None:
                filename = item.get("file") or item.get("path")
                if not filename:
                    raise AudioError(f"{label_file.name}:{line_no}: brak features/file")
                path = _safe_child(data_dir, str(filename))
                if not path.exists():
                    raise AudioError(f"{label_file.name}:{line_no}: brak pliku {path.name}")
                if path.suffix.lower() != ".wav":
                    raise AudioError(f"{label_file.name}:{line_no}: dozwolone są pliki WAV")
                features = wav_features(path, input_dim)
            if not isinstance(features, list) or len(features) != input_dim:
                raise AudioError(f"{label_file.name}:{line_no}: zły wymiar cech")
            values = []
            for v in features:
                f = float(v)
                if not math.isfinite(f):
                    raise AudioError(f"{label_file.name}:{line_no}: NaN/Inf w cechach")
                values.append(f)
            rows.append({
                "features": values, "label": label,
                "split": str(item.get("split") or "train"),
                "group": str(item.get("group") or label_file.stem),
            })
            if len(rows) > MAX_RECORDS:
                raise AudioError("Zbyt wiele rekordów audio.")
    if not rows:
        raise AudioError("Brak rekordów audio.")
    return rows
