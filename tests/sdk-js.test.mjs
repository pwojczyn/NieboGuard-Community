import test from 'node:test';
import assert from 'node:assert/strict';
import { buildUrl, NieboGuardClient } from '../sdk/javascript/client.mjs';

test('buildUrl skips empty values', () => {
  const url = buildUrl('https://api.nieboguard.pl', '/api/v2/feed', { limit: 5, region: '', kind: null });
  assert.equal(url, 'https://api.nieboguard.pl/api/v2/feed?limit=5');
});

test('client is read-only and surfaces errors', async () => {
  const calls = [];
  const fakeFetch = async (url) => {
    calls.push(url);
    return new Response(JSON.stringify({ error: { code: 'rate_limited', message: 'slow down' } }), { status: 429 });
  };
  const client = new NieboGuardClient({ fetchImpl: fakeFetch });
  assert.equal(typeof client.feed, 'function');
  await assert.rejects(() => client.feed({ limit: 1 }), (error) => error.status === 429 && error.code === 'rate_limited');
  assert.equal(calls.length, 1);
  assert.match(calls[0], /\/api\/v2\/feed\?limit=1$/);
});
