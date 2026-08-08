package cache

import (
	"testing"

	"github.com/Dimensionexpert/movieLibrary/internal/movie"
)

func TestGet_MissOnEmptyCache(t *testing.T) {
	c := NewCache(2)

	_, found := c.Get(1)
	if found {
		t.Error("expected miss on empty cache, got found=true")
	}
}

func TestPutThenGet_Hit(t *testing.T) {
	c := NewCache(2)
	m := movie.Movie{ID: 1, Title: "Inception"}

	c.Put(1, m)

	got, found := c.Get(1)
	if !found {
		t.Fatal("expected hit after Put, got miss")
	}
	if got.Title != "Inception" {
		t.Errorf("expected title 'Inception', got %q", got.Title)
	}
}

func TestPut_EvictsLeastRecentlyUsed(t *testing.T) {
	c := NewCache(2)
	c.Put(1, movie.Movie{ID: 1, Title: "A"})
	c.Put(2, movie.Movie{ID: 2, Title: "B"})
	// cache full: [2, 1] (2 = most recent)

	c.Put(3, movie.Movie{ID: 3, Title: "C"})
	// adding 3 should evict 1 (least recently used, never touched since insert)

	if _, found := c.Get(1); found {
		t.Error("expected key 1 to be evicted, but it was found")
	}
	if _, found := c.Get(2); !found {
		t.Error("expected key 2 to still be present")
	}
	if _, found := c.Get(3); !found {
		t.Error("expected key 3 to be present")
	}

	// map should not have a stale entry for the evicted key
	if _, exists := c.item[1]; exists {
		t.Error("expected evicted key 1 to be removed from the map, but it still exists")
	}
}

func TestGet_PromotesRecency(t *testing.T) {
	c := NewCache(3)
	c.Put(1, movie.Movie{ID: 1, Title: "A"})
	c.Put(2, movie.Movie{ID: 2, Title: "B"})
	c.Put(3, movie.Movie{ID: 3, Title: "C"})
	// cache full: [3, 2, 1] (1 = least recently used)

	// touch 1 via Get — this should move it to front, saving it from eviction
	if _, found := c.Get(1); !found {
		t.Fatal("expected key 1 to be found before eviction test")
	}

	c.Put(4, movie.Movie{ID: 4, Title: "D"})
	// now 2 should be the least recently used (never touched after insert), not 1

	if _, found := c.Get(2); found {
		t.Error("expected key 2 to be evicted (least recently used), but it was found")
	}
	if _, found := c.Get(1); !found {
		t.Error("expected key 1 to still be present (was touched via Get, so protected)")
	}
	if _, found := c.Get(4); !found {
		t.Error("expected key 4 to be present")
	}
}

func TestPut_UpdateExistingKey_NoDuplicate(t *testing.T) {
	c := NewCache(2)
	c.Put(1, movie.Movie{ID: 1, Title: "Old Title"})
	c.Put(1, movie.Movie{ID: 1, Title: "New Title"})

	got, found := c.Get(1)
	if !found {
		t.Fatal("expected key 1 to be found")
	}
	if got.Title != "New Title" {
		t.Errorf("expected updated title 'New Title', got %q", got.Title)
	}

	// the whole point of the bug we fixed: updating an existing key should
	// NOT create a second node or trigger eviction of anything else.
	if c.order.Len() != 1 {
		t.Errorf("expected list length 1 after updating existing key, got %d", c.order.Len())
	}
	if len(c.item) != 1 {
		t.Errorf("expected map length 1 after updating existing key, got %d", len(c.item))
	}
}
