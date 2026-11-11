/**
 * NieboGuard Federation client — ngf/1 transport (public, no storage).
 *
 * This SDK deliberately does NOT persist tokens, observations or outbox data.
 * The host application must provide secure storage and a bounded outbox.
 * Never put refresh tokens in localStorage, logs, URLs or committed files.
 */
const DEFAULT_BASE = 'https://api.nieboguard.pl';

export class FederationError extends Error {
  constructor(message, status, code) {
    super(message);
    this.name = 'FederationError';
    this.status = status;
    this.code = code || null;
  }
}

function normalizeBase(raw) {
  const value = String(raw || DEFAULT_BASE).trim();
  const url = new URL(value);
  const loopback = ['127.0.0.1', 'localhost', '::1', '[::1]'].includes(url.hostname);
  if (url.protocol !== 'https:' && !(url.protocol === 'http:' && loopback)) {
    throw new Error('federation base URL must use HTTPS');
  }
  if (url.username || url.password) throw new Error('federation base URL cannot contain credentials');
  return url.origin + url.pathname.replace(/\/+$/, '');
}
function requestId() {
  const bytes = new Uint8Array(12);
  globalThis.crypto.getRandomValues(bytes);
  return `req_${[...bytes].map((b) => b.toString(16).padStart(2, '0')).join('')}`;
}

export class FederationClient {
  constructor(options = {}) {
    this.base = normalizeBase(options.baseUrl);
    this.fetchImpl = options.fetchImpl || globalThis.fetch;
    this.timeoutMs = Math.max(1000, Math.min(60000, Number(options.timeoutMs) || 12000));
    this.nodeId = options.nodeId || null;
    this.accessToken = options.accessToken || null;
    this.refreshToken = options.refreshToken || null;
    if (typeof this.fetchImpl !== 'function') throw new Error('fetch is required');
  }

  status() {
    return {
      configured: !!this.base,
      node_id: this.nodeId,
      has_access_token: !!this.accessToken,
      has_refresh_token: !!this.refreshToken,
      protocol: 'ngf/1',
    };
  }

  credentials() {
    return { node_id: this.nodeId, access_token: this.accessToken, refresh_token: this.refreshToken };
  }

  setCredentials(value = {}) {
    if ('nodeId' in value) this.nodeId = value.nodeId || null;
    if ('accessToken' in value) this.accessToken = value.accessToken || null;
    if ('refreshToken' in value) this.refreshToken = value.refreshToken || null;
  }

  async request(path, body = null, options = {}) {
    const auth = options.auth !== false;
    const retry = options.retry !== false;
    if (auth && !this.accessToken && this.refreshToken) await this.refresh();
    const headers = { Accept: 'application/json', 'X-NG-Protocol': 'ngf/1', 'X-NG-Request-Id': requestId() };
    if (auth && this.accessToken) headers.Authorization = `Bearer ${this.accessToken}`;
    if (body != null) headers['Content-Type'] = 'application/json';
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), this.timeoutMs);
    let response;
    try {
      response = await this.fetchImpl(this.base + path, {
        method: body == null && options.method === 'GET' ? 'GET' : 'POST',
        headers, body: body == null ? undefined : JSON.stringify(body), signal: controller.signal,
      });
    } finally {
      clearTimeout(timer);
    }
    const text = await response.text();
    let data = null;
    try { data = text ? JSON.parse(text) : null; } catch { data = null; }
    if (response.status === 401 && auth && retry && this.refreshToken) {
      await this.refresh();
      return this.request(path, body, { ...options, retry: false });
    }
    if (!response.ok) {
      throw new FederationError(data?.error?.message || `HTTP ${response.status}`, response.status, data?.error?.code);
    }
    return data;
  }

  async enroll(input = {}) {
    if (!input.inviteCode) throw new Error('inviteCode is required');
    const data = await this.request('/api/v2/federation/nodes/enroll', {
      invite_code: String(input.inviteCode),
      name: String(input.name || 'NieboGuard Node').slice(0, 120),
      node_type: String(input.nodeType || 'linux'),
      version: String(input.version || '0.1.0').slice(0, 40),
      os: input.os ? String(input.os).slice(0, 40) : undefined,
      arch: input.arch ? String(input.arch).slice(0, 40) : undefined,
      capabilities: input.capabilities && typeof input.capabilities === 'object' ? input.capabilities : {},
    }, { auth: false });
    this.nodeId = data?.node_id || null;
    this.accessToken = data?.access_token || null;
    this.refreshToken = data?.refresh_token || null;
    return data;
  }

  async refresh() {
    if (!this.nodeId || !this.refreshToken) throw new Error('nodeId and refreshToken are required');
    const data = await this.request('/api/v2/federation/nodes/token', {
      node_id: this.nodeId, refresh_token: this.refreshToken,
    }, { auth: false, retry: false });
    this.accessToken = data?.access_token || null;
    this.refreshToken = data?.refresh_token || this.refreshToken;
    return data;
  }

  heartbeat(metrics = {}, version = '0.1.0') {
    return this.request('/api/v2/federation/nodes/heartbeat', { version, metrics });
  }

  config() {
    return this.request('/api/v2/federation/nodes/config', null, { method: 'GET' });
  }

  syncPull() {
    return this.request('/api/v2/federation/sync/pull', {});
  }

  syncAck(batchId, cursor = '') {
    if (!batchId) throw new Error('batchId is required');
    return this.request('/api/v2/federation/sync/ack', { batch_id: String(batchId), cursor: String(cursor || '') });
  }

  ingest(batchId, observations) {
    if (!batchId) throw new Error('batchId is required');
    if (!Array.isArray(observations) || !observations.length) throw new Error('observations must be a non-empty array');
    if (observations.length > 200) throw new Error('at most 200 observations per batch');
    return this.request('/api/v2/federation/ingest/observations', {
      batch_id: String(batchId), observations,
    });
  }

  connectorStatus(connectorId, status, error = '') {
    if (!connectorId) throw new Error('connectorId is required');
    return this.request('/api/v2/federation/nodes/connectors/status', {
      connector_id: String(connectorId), status: String(status), error: String(error || '').slice(0, 300),
    });
  }

  async revoke() {
    const data = await this.request('/api/v2/federation/nodes/revoke', {});
    this.accessToken = null;
    this.refreshToken = null;
    return data;
  }
}

export default FederationClient;
