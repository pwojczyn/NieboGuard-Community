"""Minimalny audyt sprzętu potrzebny do przydziału zadań.

Nie zbieramy nazwy hosta, adresów MAC, numerów seryjnych, ścieżek użytkownika
ani dokładnej lokalizacji. RAM i VRAM są zaokrąglane do kubełków.
"""
from __future__ import annotations

import ctypes
import os
import platform
import shutil
import subprocess
import sys
from pathlib import Path
from typing import Any


def _run(cmd: list[str], timeout: float = 8.0) -> str:
    try:
        out = subprocess.run(cmd, capture_output=True, text=True, timeout=timeout, check=False)
        return (out.stdout or "").strip()
    except Exception:
        return ""


def _ram_gb() -> float:
    try:
        if os.name == "nt":
            class MEMORYSTATUSEX(ctypes.Structure):
                _fields_ = [
                    ("dwLength", ctypes.c_ulong),
                    ("dwMemoryLoad", ctypes.c_ulong),
                    ("ullTotalPhys", ctypes.c_ulonglong),
                    ("ullAvailPhys", ctypes.c_ulonglong),
                    ("ullTotalPageFile", ctypes.c_ulonglong),
                    ("ullAvailPageFile", ctypes.c_ulonglong),
                    ("ullTotalVirtual", ctypes.c_ulonglong),
                    ("ullAvailVirtual", ctypes.c_ulonglong),
                    ("ullAvailExtendedVirtual", ctypes.c_ulonglong),
                ]
            stat = MEMORYSTATUSEX()
            stat.dwLength = ctypes.sizeof(MEMORYSTATUSEX)
            ctypes.windll.kernel32.GlobalMemoryStatusEx(ctypes.byref(stat))
            return stat.ullTotalPhys / (1024 ** 3)
        if sys.platform.startswith("linux"):
            for line in Path("/proc/meminfo").read_text(encoding="utf-8").splitlines():
                if line.startswith("MemTotal:"):
                    return int(line.split()[1]) / (1024 ** 2)
        return (os.sysconf("SC_PAGE_SIZE") * os.sysconf("SC_PHYS_PAGES")) / (1024 ** 3)
    except Exception:
        return 0.0


def _bucket(value: float, buckets: list[int]) -> int:
    out = 0
    for b in buckets:
        if value >= b:
            out = b
    return out


def _nvidia() -> dict[str, Any] | None:
    if not shutil.which("nvidia-smi"):
        return None
    raw = _run(["nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader,nounits"])
    if not raw:
        return None
    first = raw.splitlines()[0]
    parts = [x.strip() for x in first.split(",")]
    if len(parts) < 2:
        return None
    try:
        memory = float(parts[1]) / 1024.0
    except ValueError:
        memory = 0.0
    return {
        "vendor": "nvidia",
        "model": parts[0][:80],
        "memory_gb": _bucket(memory, [2, 4, 6, 8, 12, 16, 24, 48, 80]),
    }


def _other_gpu() -> dict[str, Any] | None:
    try:
        if os.name == "nt":
            raw = _run(["powershell", "-NoProfile", "-Command",
                        "Get-CimInstance Win32_VideoController | Select-Object -ExpandProperty Name"])
        elif sys.platform.startswith("linux") and shutil.which("lspci"):
            raw = _run(["lspci"])
            raw = "\n".join(line.split(": ", 1)[1] for line in raw.splitlines() if "VGA" in line or "3D controller" in line)
        else:
            raw = ""
        name = raw.splitlines()[0].strip() if raw else ""
        if not name:
            return None
        low = name.lower()
        if "nvidia" in low:
            return None
        return {"vendor": "other", "model": name[:80], "memory_gb": None}
    except Exception:
        return None


def detect() -> dict[str, Any]:
    system = platform.system().lower()
    os_family = "windows" if system == "windows" else "linux" if system == "linux" else system
    arch = platform.machine().lower()
    if arch in ("x86_64", "amd64"):
        arch = "x64"
    elif arch in ("aarch64", "arm64"):
        arch = "arm64"
    ram = _ram_gb()
    gpu = _nvidia() or _other_gpu()
    return {
        "os_family": os_family,
        "os_version": platform.release()[:80],
        "arch": arch,
        "cpu_cores": os.cpu_count() or 1,
        "ram_gb": _bucket(ram, [4, 8, 16, 32, 64, 128, 256, 512, 1024]),
        "ram_gb_actual": round(ram, 1),
        "gpu_vendor": (gpu or {}).get("vendor", "none"),
        "gpu_model": (gpu or {}).get("model"),
        "gpu_memory_gb": (gpu or {}).get("memory_gb"),
        "python": platform.python_version(),
        "supported": os_family in ("windows", "linux") and arch in ("x64", "arm64"),
    }


def idle_seconds() -> float | None:
    """Return desktop idle time, or None when the platform cannot report it.

    A None value must not be treated as "idle"; the runner refuses to start in
    idle mode without a real signal and asks for --schedule manual/always.
    """
    try:
        if os.name == "nt":
            class LASTINPUTINFO(ctypes.Structure):
                _fields_ = [("cbSize", ctypes.c_uint), ("dwTime", ctypes.c_uint)]
            info = LASTINPUTINFO()
            info.cbSize = ctypes.sizeof(LASTINPUTINFO)
            if not ctypes.windll.user32.GetLastInputInfo(ctypes.byref(info)):
                return None
            millis = ctypes.windll.kernel32.GetTickCount() - info.dwTime
            return max(0.0, millis / 1000.0)
        if sys.platform.startswith("linux"):
            xidle = shutil.which("xprintidle")
            if not xidle:
                return None
            raw = _run([xidle], timeout=2)
            return max(0.0, int(raw) / 1000.0) if raw.isdigit() else None
    except Exception:
        return None
    return None


def summary(hw: dict[str, Any]) -> str:
    gpu = "brak"
    if hw.get("gpu_vendor") == "nvidia":
        gpu = f"NVIDIA {hw.get('gpu_model') or ''} ({hw.get('gpu_memory_gb') or '?'} GB)"
    elif hw.get("gpu_vendor") == "other":
        gpu = f"{hw.get('gpu_model') or 'inne GPU'} — niepotwierdzone"
    return (f"{hw.get('os_family')} {hw.get('arch')} · CPU {hw.get('cpu_cores')} rdzeni · "
            f"RAM {hw.get('ram_gb')} GB · GPU: {gpu} · Python {hw.get('python')}")
