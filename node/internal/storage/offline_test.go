package storage

import (
	"testing"
	"time"
)

func TestOfflinePackLifecycle(t *testing.T) {
	store, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := store.Initialize(); err != nil {
		t.Fatal(err)
	}
	pack := &OfflinePack{ID: "pack-a", Version: "1.0.0", Sha256: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Bytes: 10, CreatedAt: "2026-11-01T00:00:00Z", ManifestJSON: "{}", Mode: "REPLAY", Sharing: "LOCAL_ONLY", Status: "AVAILABLE", UpdatedAt: "2026-11-01T00:00:00Z"}
	if err := store.UpsertOfflinePack(pack); err != nil {
		t.Fatal(err)
	}
	pack.Version = "1.1.0"
	pack.CreatedAt = "2026-11-02T00:00:00Z"
	if err := store.UpsertOfflinePack(pack); err != nil {
		t.Fatal(err)
	}
	packs, err := store.ListOfflinePacks("pack-a", 10)
	if err != nil || len(packs) != 2 {
		t.Fatalf("list: %v %d", err, len(packs))
	}
	current, err := store.CurrentOfflinePack("pack-a")
	if err != nil || current == nil || current.Version != "1.1.0" {
		t.Fatalf("current: %v %#v", err, current)
	}
	n, err := store.PruneOfflinePacks(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil || n != 2 {
		t.Fatalf("prune: %v %d", err, n)
	}
	if packs, _ := store.ListOfflinePacks("", 10); len(packs) != 0 {
		t.Fatalf("not pruned: %d", len(packs))
	}
}
