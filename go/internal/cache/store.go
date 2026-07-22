package cache

import (
	"errors"
	"sync"

	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/profile"
)

// Store is a concurrency-safe in-memory profile cache.
type Store struct {
	mu        sync.RWMutex
	profiles  map[string]*profile.Profile
	recentIDs []string
}

// NewStore returns an initialized profile store.
func NewStore() *Store {
	return &Store{profiles: make(map[string]*profile.Profile)}
}

// Put stores a defensive copy of a profile. The zero value is also usable.
func (s *Store) Put(value profile.Profile) error {
	if value.ID == "" {
		return errors.New("profile ID is required")
	}

	copyValue := cloneProfile(value)
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.profiles[value.ID]; !exists {
		s.recentIDs = append(s.recentIDs, value.ID)
	}
	s.profiles[value.ID] = &copyValue
	return nil
}

// Get returns a defensive copy and reports whether a non-nil entry exists.
func (s *Store) Get(id string) (profile.Profile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.profiles[id]
	if !exists {
		return profile.Profile{}, false
	}
	return cloneProfile(*value), true
}

// MostRecent returns the newest stored profile ID when one exists.
func (s *Store) MostRecent() (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.recentIDs[len(s.recentIDs)-1], true
}

// RecentIDs returns a copy that callers may mutate safely.
func (s *Store) RecentIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return append([]string(nil), s.recentIDs...)
}

func cloneProfile(value profile.Profile) profile.Profile {
	value.Tags = append([]string(nil), value.Tags...)
	return value
}
