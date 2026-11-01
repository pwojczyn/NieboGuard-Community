# Sample data

All files here are **synthetic**. They are not real flights, alerts, users or
locations. Every top-level object includes `"synthetic": true`.

- `adsb-aircraft.sample.json` — dump1090/readsb-like input; includes invalid
  rows to test filtering.
- `feed.sample.json` — public-feed-like items.
- `state.sample.json` — minimal map snapshot.

Do not replace these files with real user data or private receiver locations.
If you add a real dataset, include source, license and attribution in a PR.
