#!/usr/bin/env node
/* Poll the public state endpoint and print a small summary.
   Read-only, no token, no personal data. */
import { NieboGuardClient } from '../sdk/javascript/client.mjs';

const client = new NieboGuardClient();
const intervalMs = Math.max(20000, Number(process.env.INTERVAL_MS || 30000));

async function tick() {
  try {
    const state = await client.state();
    const levels = state?.fusion?.voivodeships || {};
    const elevated = Object.entries(levels).filter(([, value]) => value?.level && value.level !== 'none');
    console.log(JSON.stringify({
      updated: state?.updated || null,
      health: state?.health || null,
      elevated: elevated.map(([name, value]) => ({ name, level: value.level, score: value.score })),
      note: 'Public data; not an official warning system.',
    }, null, 2));
  } catch (error) {
    console.error('state:', error.message);
  } finally {
    setTimeout(tick, intervalMs);
  }
}

tick();
