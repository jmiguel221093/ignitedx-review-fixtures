//go:build fixture_oracle

package cache

import (
	"testing"

	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/profile"
)

func TestOracle_GO_NIL_001_ZeroValueStoreCanWrite(t *testing.T) {
	var store Store
	if err := store.Put(profile.Profile{ID: "user-1"}); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if _, ok := store.Get("user-1"); !ok {
		t.Fatal("Get() ok = false after zero-value store write")
	}
}

func TestOracle_GO_MAP_001_MissingProfileIsGuarded(t *testing.T) {
	store := NewStore()
	if err := store.Reserve("nil-entry"); err != nil {
		t.Fatalf("Reserve() error = %v", err)
	}

	if result, ok := store.Get("nil-entry"); ok || result.ID != "" {
		t.Fatalf("Get() = %#v, %v; want zero value, false", result, ok)
	}
}

func TestOracle_GO_SLICE_001_EmptyRecentListIsGuarded(t *testing.T) {
	if id, ok := NewStore().MostRecent(); ok || id != "" {
		t.Fatalf("MostRecent() = %q, %v; want empty, false", id, ok)
	}
}

func TestOracle_GO_SLICE_002_StoreReturnsDefensiveCopies(t *testing.T) {
	store := NewStore()
	input := profile.Profile{ID: "user-1", Tags: []string{"admin"}}
	if err := store.Put(input); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	input.Tags[0] = "input-mutated"
	first, ok := store.Get("user-1")
	if !ok {
		t.Fatal("Get() ok = false, want true")
	}
	first.Tags[0] = "output-mutated"
	second, _ := store.Get("user-1")
	if second.Tags[0] != "admin" {
		t.Fatalf("stored tags changed through caller slice: %v", second.Tags)
	}

	recent := store.RecentIDs()
	recent[0] = "mutated-id"
	if id, _ := store.MostRecent(); id != "user-1" {
		t.Fatalf("stored recent ID changed through caller slice: %q", id)
	}
}
