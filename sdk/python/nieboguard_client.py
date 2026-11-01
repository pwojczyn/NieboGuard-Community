"""Read-only Python client for the public NieboGuard API.

No API key, no writes, no user data. Standard library only.
"""
from __future__ import annotations

import json
import urllib.error
import urllib.parse
import urllib.request
from typing import Any

DEFAULT_BASE = "https://api.nieboguard.pl"


class NieboGuardError(RuntimeError):
    def __init__(self, message: str, status: int | None = None, code: str | None = None):
        super().__init__(message)
        self.status = status
        self.code = code


def build_url(base: str, path: str, params: dict[str, Any] | None = None) -> str:
    root = (base or DEFAULT_BASE).rstrip("/") + "/"
    url = urllib.parse.urljoin(root, path.lstrip("/"))
    clean = {k: v for k, v in (params or {}).items() if v is not None and v != ""}
    if clean:
        url += ("&" if "?" in url else "?") + urllib.parse.urlencode(clean)
    return url


class NieboGuardClient:
    def __init__(self, base_url: str = DEFAULT_BASE, timeout: float = 12.0):
        self.base_url = (base_url or DEFAULT_BASE).rstrip("/")
        self.timeout = timeout

    def request(self, path: str, params: dict[str, Any] | None = None) -> Any:
        req = urllib.request.Request(
            build_url(self.base_url, path, params),
            headers={"Accept": "application/json", "User-Agent": "NieboGuard-Community-SDK/0.1"},
        )
        try:
            with urllib.request.urlopen(req, timeout=self.timeout) as response:
                body = response.read().decode("utf-8")
                return json.loads(body) if body else None
        except urllib.error.HTTPError as exc:
            try:
                data = json.loads(exc.read().decode("utf-8"))
            except Exception:
                data = {}
            raise NieboGuardError(
                (data.get("error") or {}).get("message") or f"HTTP {exc.code}",
                exc.code,
                (data.get("error") or {}).get("code"),
            ) from exc
        except urllib.error.URLError as exc:
            raise NieboGuardError(f"Network error: {exc.reason}") from exc

    def state(self) -> Any:
        return self.request("/api/state")

    def feed(self, **params: Any) -> Any:
        return self.request("/api/v2/feed", params)

    def alarms(self, **params: Any) -> Any:
        return self.request("/api/v2/alarms", params)

    def incidents(self, **params: Any) -> Any:
        return self.request("/api/v2/incidents", params)

    def stats(self) -> Any:
        return self.request("/api/v2/stats/overview")

    def sources(self) -> Any:
        return self.request("/api/v2/sources")

    def source(self, source_id: str) -> Any:
        return self.request(f"/api/v2/sources/{urllib.parse.quote(str(source_id), safe='')}")

    def locations(self) -> Any:
        return self.request("/api/v2/locations")

    def threats(self) -> Any:
        return self.request("/api/v2/threats")
