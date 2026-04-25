package storage

import "testing"

func TestStoreCRUDAndExists(t *testing.T) {
	store := NewStore()

	if _, ok := store.Get("missing"); ok {
		t.Fatal("expected missing key to be absent")
	}

	store.Set("k1", "v1")
	store.Set("k2", "v2")

	if got, ok := store.Get("k1"); !ok || got != "v1" {
		t.Fatalf("expected k1=v1, got %q present=%v", got, ok)
	}

	if got := store.Exists("k1", "k2", "k3"); got != 2 {
		t.Fatalf("expected exists count 2, got %d", got)
	}

	if got := store.Del("k2", "k3"); got != 1 {
		t.Fatalf("expected del count 1, got %d", got)
	}

	if got := store.Exists("k2"); got != 0 {
		t.Fatalf("expected k2 deleted, exists=%d", got)
	}
}