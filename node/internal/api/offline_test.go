package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pwojczyn/NieboGuard-Community/node/internal/storage"
)

func offlineStore(t *testing.T) *storage.Storage {
	store, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Initialize(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestOfflinePackAPI(t *testing.T) {
	s := &Server{storage: offlineStore(t)}
	body := `{"id":"pack-a","version":"1.0.0","sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","bytes":10,"mode":"REPLAY","sharing":"LOCAL_ONLY"}`
	r := httptest.NewRequest(http.MethodPost, "/api/v1/admin/offline/packs", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	s.AdminOfflinePackHandler(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status %d: %s", w.Code, w.Body.String())
	}
	r = httptest.NewRequest(http.MethodGet, "/api/v1/offline/packs?id=pack-a", nil)
	w = httptest.NewRecorder()
	s.OfflinePacksHandler(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("list status %d", w.Code)
	}
	var listed struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &listed); err != nil || listed.Total != 1 {
		t.Fatalf("list payload %s", w.Body.String())
	}
	r = httptest.NewRequest(http.MethodGet, "/api/v1/offline/packs/pack-a", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "pack-a"})
	w = httptest.NewRecorder()
	s.CurrentOfflinePackHandler(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("current status %d", w.Code)
	}
}

func TestOfflinePackRejectsInvalidManifest(t *testing.T) {
	s := &Server{storage: offlineStore(t)}
	for _, body := range []string{
		`{"id":"pack","version":"1","sha256":"bad","bytes":1}`,
		`{"id":"pack","version":"1","sha256":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","bytes":1,"mode":"SIMULATED","sharing":"COMMUNITY_SHARED"}`,
	} {
		r := httptest.NewRequest(http.MethodPost, "/api/v1/admin/offline/packs", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		s.AdminOfflinePackHandler(w, r)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d for %s", w.Code, body)
		}
	}
}
