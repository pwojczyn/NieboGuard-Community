# Examples

- [`state_watch.mjs`](state_watch.mjs) — poll `/api/state` every 30 s and print
  a small summary.
- [`feed_to_csv.mjs`](feed_to_csv.mjs) — export public feed items to CSV.

```bash
node examples/state_watch.mjs
node examples/feed_to_csv.mjs /tmp/feed.csv
```

Examples are read-only and use no token. Do not add personal data or private
coordinates to outputs. This is not an official warning system.
