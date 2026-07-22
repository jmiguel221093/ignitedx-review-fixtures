package cache

import (
	"testing"

	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/profile"
)

func TestStorePutAndGet(t *testing.T) {
	t.Parallel()

	store := NewStore()
	input := profile.Profile{ID: "user-1", Name: "Ada", Tags: []string{"admin"}}
	if err := store.Put(input); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	result, ok := store.Get("user-1")
	if !ok {
		t.Fatal("Get() ok = false, want true")
	}
	if result.ID != input.ID || result.Name != input.Name {
		t.Fatalf("Get() result = %#v, want %#v", result, input)
	}
}

func TestStoreMostRecentReturnsStoredProfile(t *testing.T) {
	t.Parallel()

	store := NewStore()
	if err := store.Put(profile.Profile{ID: "user-1"}); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if id, ok := store.MostRecent(); !ok || id != "user-1" {
		t.Fatalf("MostRecent() = %q, %v; want user-1, true", id, ok)
	}
}
