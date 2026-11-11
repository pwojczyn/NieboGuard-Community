import test from 'node:test';
import assert from 'node:assert/strict';
import { FederationClient, FederationError } from '../sdk/javascript/federation.mjs';

function response(status, value) {
  return new Response(JSON.stringify(value), { status, headers: { 'content-type': 'application/json' } });
}

test('public federation SDK performs enrollment, refresh, heartbeat, config and revoke', async () => {
  const calls = [];
  let access = 'access-1';
  const client = new FederationClient({
    baseUrl: 'https://api.example.test',
    fetchImpl: async (url, init) => {
      const path = new URL(url).pathname;
      calls.push({ path, method: init.method, headers: init.headers, body: init.body ? JSON.parse(init.body) : null });
      if (path.endsWith('/enroll')) return response(201, { node_id: 'fnode_1', organization_id: 'forg_1', workspace_id: null, access_token: access, refresh_token: 'refresh-1', access_expires_at: '2026-11-11T10:20:00Z', protocol: 'ngf/1' });
      if (path.endsWith('/token')) return response(200, { access_token: 'access-2', refresh_token: 'refresh-2', access_expires_at: '2026-11-11T10:20:00Z' });
      if (path.endsWith('/heartbeat')) {
        if (access === 'access-1') { access = 'force-refresh'; return response(401, { error: { code: 'expired' } }); }
        return response(200, { ok: true, config_version: 1 });
      }
      if (path.endsWith('/config')) {
        assert.equal(init.method, 'GET');
        return response(200, { config_version: 1, config: {} });
      }
      if (path.endsWith('/revoke')) return response(200, { ok: true, revoked: true });
      return response(404, { error: { code: 'not_found' } });
    }
  });
  assert.equal(client.status().node_id, null);
  await client.enroll({ inviteCode: 'a'.repeat(64), name: 'Public SDK' });
  assert.equal(client.status().node_id, 'fnode_1');
  assert.equal(client.credentials().refresh_token, 'refresh-1');
  assert.equal(calls[0].headers['X-NG-Protocol'], 'ngf/1');
  assert.match(calls[0].headers['X-NG-Request-Id'], /^req_[0-9a-f]+$/);

  await client.heartbeat({ queue: 0 }, '0.1.0');
  assert(calls.some((c) => c.path.endsWith('/token')), '401 triggers refresh');
  assert.equal(client.credentials().access_token, 'access-2');

  await client.config();
  await client.revoke();
  assert.equal(client.status().has_refresh_token, false);
  assert.equal(client.status().has_access_token, false);
});

test('public federation SDK validates requests before network access', async () => {
  let called = 0;
  const client = new FederationClient({
    baseUrl: 'https://api.example.test',
    fetchImpl: async () => { called++; return response(200, {}); }
  });
  assert.throws(() => client.ingest('b1', []), /observations/);
  assert.throws(() => client.ingest('b1', Array.from({ length: 201 }, (_, i) => ({ event_id: String(i) }))), /200/);
  assert.throws(() => client.syncAck(''), /batchId/);
  assert.equal(called, 0);
  assert.throws(() => new FederationClient({ baseUrl: 'http://example.com' }), /HTTPS/);
  assert.doesNotThrow(() => new FederationClient({ baseUrl: 'http://127.0.0.1:8787' }));
  const denied = new FederationClient({ baseUrl: 'https://api.example.test', fetchImpl: async () => response(403, { error: { code: 'forbidden', message: 'denied' } }) });
  await assert.rejects(() => denied.syncPull(), (error) => error instanceof FederationError && error.status === 403 && error.code === 'forbidden');
});

test('public federation SDK does not persist tokens in the SDK', () => {
  const client = new FederationClient({ baseUrl: 'https://api.example.test', fetchImpl: async () => response(200, {}) });
  client.setCredentials({ nodeId: 'n', accessToken: 'a', refreshToken: 'r' });
  assert.deepEqual(Object.keys(client.credentials()).sort(), ['access_token', 'node_id', 'refresh_token']);
  const status = client.status();
  assert.deepEqual(status, { configured: true, node_id: 'n', has_access_token: true, has_refresh_token: true, protocol: 'ngf/1' });
  assert.equal(JSON.stringify(status).includes('"r"'), false);
  assert.equal(JSON.stringify(status).includes('"a"'), false);
});
