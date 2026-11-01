/**
 * Small read-only client for the public NieboGuard API.
 * No API key, no write operations and no user data.
 */
const DEFAULT_BASE = 'https://api.nieboguard.pl';

export function buildUrl(base, path, params = {}) {
  const url = new URL(path, String(base || DEFAULT_BASE).replace(/\/+$/, '') + '/');
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null || value === '') continue;
    url.searchParams.set(key, String(value));
  }
  return url.toString();
}

export class NieboGuardError extends Error {
  constructor(message, status, code) {
    super(message);
    this.name = 'NieboGuardError';
    this.status = status;
    this.code = code || null;
  }
}

export class NieboGuardClient {
  constructor(options = {}) {
    this.base = String(options.baseUrl || DEFAULT_BASE).replace(/\/+$/, '');
    this.timeoutMs = Number(options.timeoutMs || 12000);
    this.fetchImpl = options.fetchImpl || globalThis.fetch;
    if (typeof this.fetchImpl !== 'function') throw new Error('fetch is required');
  }

  async request(path, params = {}) {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), this.timeoutMs);
    try {
      const response = await this.fetchImpl(buildUrl(this.base, path, params), {
        method: 'GET',
        headers: { Accept: 'application/json' },
        signal: controller.signal,
      });
      const text = await response.text();
      let data = null;
      try { data = text ? JSON.parse(text) : null; } catch { data = null; }
      if (!response.ok) {
        throw new NieboGuardError(
          data?.error?.message || `HTTP ${response.status}`,
          response.status,
          data?.error?.code,
        );
      }
      return data;
    } finally {
      clearTimeout(timer);
    }
  }

  state() { return this.request('/api/state'); }
  feed(params = {}) { return this.request('/api/v2/feed', params); }
  alarms(params = {}) { return this.request('/api/v2/alarms', params); }
  incidents(params = {}) { return this.request('/api/v2/incidents', params); }
  stats() { return this.request('/api/v2/stats/overview'); }
  sources() { return this.request('/api/v2/sources'); }
  source(id) { return this.request(`/api/v2/sources/${encodeURIComponent(id)}`); }
  locations() { return this.request('/api/v2/locations'); }
  threats() { return this.request('/api/v2/threats'); }
}

export default NieboGuardClient;
