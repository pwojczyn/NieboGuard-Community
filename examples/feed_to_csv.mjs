#!/usr/bin/env node
/* Export public feed items to CSV. Read-only; no token or personal data. */
import { writeFileSync } from 'node:fs';
import { NieboGuardClient } from '../sdk/javascript/client.mjs';

const client = new NieboGuardClient();
const output = process.argv[2] || 'feed.csv';
const feed = await client.feed({ limit: 50 });
const rows = feed?.items || [];
const escape = (value) => `"${String(value ?? '').replaceAll('"', '""')}"`;
const header = ['id', 'kind', 'source', 'category', 'title', 'region', 'published_at', 'verification'];
const csv = [header.join(','), ...rows.map((row) => header.map((key) => escape(row[key])).join(','))].join('\n');
writeFileSync(output, csv + '\n', 'utf8');
console.log(`Wrote ${rows.length} rows to ${output}`);
