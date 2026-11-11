#!/usr/bin/env node
/**
 * Federation SDK public-release integrity manifest.
 *
 * This does not publish, sign or upload anything. It produces deterministic
 * SHA-256 metadata and a minimal CycloneDX SBOM for review, and fails on
 * obvious private paths/secrets in the public SDK files.
 */
import { createHash } from 'node:crypto';
import { readFileSync, statSync } from 'node:fs';
import path from 'node:path';
import { pathToFileURL } from 'node:url';

export const SDK_FILES = [
  'sdk/javascript/federation.mjs',
  'sdk/python/nieboguard_federation.py',
  'docs/FEDERATION.md',
  'docs/FEDERATION_RELEASE.md',
  'tests/federation-sdk.test.mjs',
  'tests/test_federation_sdk_python.py',
];

const FORBIDDEN = [
  /-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----/,
  /\/home\/ubuntu\//,
  /\bghp_[A-Za-z0-9]{20,}\b/,
  /\bgithub_pat_[A-Za-z0-9_]{20,}\b/,
  /\bAKIA[0-9A-Z]{16}\b/,
  /\bCF_API_TOKEN\b/,
  /\bCLOUDFLARE_API_TOKEN\b/,
  /\bOPENROUTER_API_KEY\b/,
  /\bkeystore\.properties\b/,
];

export function sha256(text) {
  return createHash('sha256').update(text, 'utf8').digest('hex');
}

export function assertSafe(relativePath, text) {
  for (const pattern of FORBIDDEN) {
    if (pattern.test(text)) throw new Error(`forbidden content in ${relativePath}: ${pattern}`);
  }
  if (/ngf\/1|ngf:1/.test(relativePath)) throw new Error('invalid path marker');
}

export function inspect(root = process.cwd()) {
  const files = [];
  for (const relative of SDK_FILES) {
    const absolute = path.join(root, relative);
    const text = readFileSync(absolute, 'utf8');
    assertSafe(relative, text);
    files.push({
      path: relative,
      bytes: statSync(absolute).size,
      sha256: sha256(text),
    });
  }
  const js = readFileSync(path.join(root, 'sdk/javascript/federation.mjs'), 'utf8');
  const py = readFileSync(path.join(root, 'sdk/python/nieboguard_federation.py'), 'utf8');
  const protocol = js.includes('ngf/1') && py.includes('ngf/1');
  if (!protocol) throw new Error('ngf/1 marker missing from a public SDK');
  return {
    schema: 'ng-federation-sdk-manifest/v1',
    version: '0.1.0-prep',
    protocol: 'ngf/1',
    release_authorized: false,
    files,
  };
}

export function createSbom(manifest) {
  return {
    bomFormat: 'CycloneDX',
    specVersion: '1.5',
    version: 1,
    metadata: {
      component: {
        type: 'library',
        name: 'nieboguard-federation-sdk',
        version: manifest.version,
      },
      properties: [
        { name: 'ngf:protocol', value: manifest.protocol },
        { name: 'ngf:release_authorized', value: String(manifest.release_authorized) },
      ],
    },
    components: manifest.files.map((file) => ({
      type: 'file',
      name: file.path,
      hashes: [{ alg: 'SHA-256', content: file.sha256 }],
    })),
  };
}

async function main() {
  const args = new Set(process.argv.slice(2));
  try {
    const manifest = inspect();
    if (args.has('--sbom')) {
      console.log(JSON.stringify(createSbom(manifest), null, 2));
    } else if (args.has('--check')) {
      console.log(`OK — federation SDK manifest: ${manifest.files.length} files, protocol ${manifest.protocol}, release_authorized=false`);
    } else {
      console.log(JSON.stringify(manifest, null, 2));
    }
  } catch (error) {
    console.error(`FAIL — ${error.message}`);
    process.exitCode = 1;
  }
}

if (process.argv[1] && import.meta.url === pathToFileURL(path.resolve(process.argv[1])).href) {
  await main();
}
