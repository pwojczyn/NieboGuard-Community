import test from 'node:test';
import assert from 'node:assert/strict';
import { inspect, createSbom, assertSafe } from '../tools/federation/release_manifest.mjs';

test('public federation SDK manifest is deterministic and not a release authorization', () => {
  const manifest = inspect();
  assert.equal(manifest.schema, 'ng-federation-sdk-manifest/v1');
  assert.equal(manifest.protocol, 'ngf/1');
  assert.equal(manifest.release_authorized, false);
  assert.equal(manifest.files.length, 6);
  for (const file of manifest.files) {
    assert.match(file.sha256, /^[0-9a-f]{64}$/);
    assert(file.bytes > 0);
  }
});

test('minimal SBOM references every public file and protocol', () => {
  const sbom = createSbom(inspect());
  assert.equal(sbom.bomFormat, 'CycloneDX');
  assert.equal(sbom.specVersion, '1.5');
  assert.equal(sbom.metadata.component.name, 'nieboguard-federation-sdk');
  assert.equal(sbom.metadata.properties.find((x) => x.name === 'ngf:protocol').value, 'ngf/1');
  assert.equal(sbom.components.length, 6);
});

test('secret and private-path checks fail closed', () => {
  assert.throws(() => assertSafe('sdk/example.mjs', 'const token = "ghp_abcdefghijklmnopqrstuvwxyz";'), /forbidden/);
  assert.throws(() => assertSafe('sdk/example.mjs', 'path=/home/ubuntu/NieboGuard'), /forbidden/);
  assert.throws(() => assertSafe('sdk/example.mjs', '-----BEGIN PRIVATE KEY-----'), /forbidden/);
  assert.doesNotThrow(() => assertSafe('sdk/example.mjs', 'const protocol = "ngf/1";'));
});
