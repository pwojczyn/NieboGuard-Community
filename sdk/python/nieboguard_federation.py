"""Public NieboGuard Federation transport client (ngf/1).

Standard library only. The SDK never persists tokens, observations or an
outbox. Host applications must use OS-protected storage and a bounded queue.
"""
from __future__ import annotations

import json
import urllib.error
import urllib.parse
import urllib.request
import uuid
from typing import Any, Callable

DEFAULT_BASE = "https://api.nieboguard.pl"


class FederationError(RuntimeError):
    def __init__(self, message: str, status: int | None = None, code: str | None = None):
        super().__init__(message)
        self.status = status
        self.code = code


def normalize_base(value: str | None) -> str:
    raw = (value or DEFAULT_BASE).strip()
    parsed = urllib.parse.urlparse(raw)
    loopback = parsed.hostname in {"127.0.0.1", "localhost", "::1"}
    if parsed.scheme != "https" and not (parsed.scheme == "http" and loopback):
        raise ValueError("federation base URL must use HTTPS")
    if parsed.username or parsed.password:
        raise ValueError("federation base URL cannot contain credentials")
    return f"{parsed.scheme}://{parsed.netloc}{parsed.path.rstrip('/')}"


class FederationClient:
    def __init__(
        self,
        base_url: str = DEFAULT_BASE,
        timeout: float = 12.0,
        opener: Callable[..., Any] | None = None,
        node_id: str | None = None,
        access_token: str | None = None,
        refresh_token: str | None = None,
    ):
        self.base_url = normalize_base(base_url)
        self.timeout = max(1.0, min(60.0, float(timeout)))
        self.opener = opener or urllib.request.urlopen
        self.node_id = node_id
        self.access_token = access_token
        self.refresh_token = refresh_token

    def status(self) -> dict[str, Any]:
        return {
            "configured": bool(self.base_url),
            "node_id": self.node_id,
            "has_access_token": bool(self.access_token),
            "has_refresh_token": bool(self.refresh_token),
            "protocol": "ngf/1",
        }

    def credentials(self) -> dict[str, str | None]:
        return {"node_id": self.node_id, "access_token": self.access_token, "refresh_token": self.refresh_token}

    def set_credentials(self, node_id: str | None = None, access_token: str | None = None, refresh_token: str | None = None) -> None:
        self.node_id = node_id
        self.access_token = access_token
        self.refresh_token = refresh_token

    def request(self, path: str, body: dict[str, Any] | None = None, *, auth: bool = True, retry: bool = True, method: str | None = None) -> Any:
        if auth and not self.access_token and self.refresh_token:
            self.refresh()
        headers = {
            "Accept": "application/json",
            "User-Agent": "NieboGuard-Community-SDK/0.1",
            "X-NG-Protocol": "ngf/1",
            "X-NG-Request-Id": f"req_{uuid.uuid4().hex[:24]}",
        }
        if auth and self.access_token:
            headers["Authorization"] = f"Bearer {self.access_token}"
        data = None
        if body is not None:
            data = json.dumps(body).encode("utf-8")
            headers["Content-Type"] = "application/json"
        request = urllib.request.Request(
            self.base_url + path, data=data, headers=headers,
            method="GET" if method == "GET" or (body is None and method is None) else "POST",
        )
        try:
            with self.opener(request, timeout=self.timeout) as response:
                raw = response.read().decode("utf-8")
                return json.loads(raw) if raw else None
        except urllib.error.HTTPError as exc:
            try:
                payload = json.loads(exc.read().decode("utf-8"))
            except Exception:
                payload = {}
            if exc.code == 401 and auth and retry and self.refresh_token:
                self.refresh()
                return self.request(path, body, auth=auth, retry=False, method=method)
            error = (payload.get("error") or {}) if isinstance(payload, dict) else {}
            raise FederationError(error.get("message") or f"HTTP {exc.code}", exc.code, error.get("code")) from exc
        except urllib.error.URLError as exc:
            raise FederationError(f"Network error: {exc.reason}") from exc

    def enroll(self, invite_code: str, name: str = "NieboGuard Node", node_type: str = "linux", version: str = "0.1.0", **metadata: Any) -> Any:
        if not invite_code:
            raise ValueError("invite_code is required")
        payload = {
            "invite_code": str(invite_code), "name": str(name)[:120], "node_type": str(node_type),
            "version": str(version)[:40],
            "capabilities": metadata.get("capabilities") if isinstance(metadata.get("capabilities"), dict) else {},
        }
        if metadata.get("os"):
            payload["os"] = str(metadata["os"])[:40]
        if metadata.get("arch"):
            payload["arch"] = str(metadata["arch"])[:40]
        result = self.request("/api/v2/federation/nodes/enroll", payload, auth=False)
        self.node_id = result.get("node_id")
        self.access_token = result.get("access_token")
        self.refresh_token = result.get("refresh_token")
        return result

    def refresh(self) -> Any:
        if not self.node_id or not self.refresh_token:
            raise ValueError("node_id and refresh_token are required")
        result = self.request(
            "/api/v2/federation/nodes/token",
            {"node_id": self.node_id, "refresh_token": self.refresh_token},
            auth=False, retry=False,
        )
        self.access_token = result.get("access_token")
        self.refresh_token = result.get("refresh_token") or self.refresh_token
        return result

    def heartbeat(self, metrics: dict[str, Any] | None = None, version: str = "0.1.0") -> Any:
        return self.request("/api/v2/federation/nodes/heartbeat", {"version": version, "metrics": metrics or {}})

    def config(self) -> Any:
        return self.request("/api/v2/federation/nodes/config", None, method="GET")

    def sync_pull(self) -> Any:
        return self.request("/api/v2/federation/sync/pull", {})

    def sync_ack(self, batch_id: str, cursor: str = "") -> Any:
        if not batch_id:
            raise ValueError("batch_id is required")
        return self.request("/api/v2/federation/sync/ack", {"batch_id": str(batch_id), "cursor": str(cursor or "")})

    def ingest(self, batch_id: str, observations: list[dict[str, Any]]) -> Any:
        if not batch_id:
            raise ValueError("batch_id is required")
        if not isinstance(observations, list) or not observations:
            raise ValueError("observations must be a non-empty list")
        if len(observations) > 200:
            raise ValueError("at most 200 observations per batch")
        return self.request("/api/v2/federation/ingest/observations", {"batch_id": str(batch_id), "observations": observations})

    def connector_status(self, connector_id: str, status: str, error: str = "") -> Any:
        if not connector_id:
            raise ValueError("connector_id is required")
        return self.request("/api/v2/federation/nodes/connectors/status", {
            "connector_id": str(connector_id), "status": str(status), "error": str(error)[:300],
        })

    def revoke(self) -> Any:
        result = self.request("/api/v2/federation/nodes/revoke", {})
        self.access_token = None
        self.refresh_token = None
        return result
