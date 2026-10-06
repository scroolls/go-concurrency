package lrucache

import (
	"fmt"
	"sync"
	"testing"
)

func TestLRU_NonPositiveCapacity(t *testing.T) {
	for _, capacity := range []int{0, -1} {
		c := New[int](capacity)
		c.Set("a", 1)
		if _, ok := c.Get("a"); ok {
			t.Fatalf("capacity %d must not retain items", capacity)
		}
		c.Clear()
	}
}

func TestLRU_EvictLeastUsed(t *testing.T) {
	c := New[int](3)

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)

	c.Get("a")
	c.Get("c")

	c.Set("d", 4)

	if _, ok := c.Get("b"); ok {
		t.Error("Expected 'b' to be evicted, but it was found in cache")
	}

	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Errorf("Expected 'a'=1, got ok=%v, v=%v", ok, v)
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Errorf("Expected 'c'=3, got ok=%v, v=%v", ok, v)
	}
	if v, ok := c.Get("d"); !ok || v != 4 {
		t.Errorf("Expected 'd'=4, got ok=%v, v=%v", ok, v)
	}
}

func TestLRU_CapacityEviction(t *testing.T) {
	c := New[string](2)

	c.Set("x", "foo")
	c.Set("y", "bar")

	c.Set("z", "baz")

	if _, ok := c.Get("x"); ok {
		t.Error("Expected 'x' to be evicted after capacity exceeded")
	}
	if v, ok := c.Get("y"); !ok || v != "bar" {
		t.Errorf("Expected 'y'='bar', got ok=%v, v=%v", ok, v)
	}
	if v, ok := c.Get("z"); !ok || v != "baz" {
		t.Errorf("Expected 'z'='baz', got ok=%v, v=%v", ok, v)
	}
}

func TestLRU_UpdateMovesToFront(t *testing.T) {
	c := New[int](2)

	c.Set("a", 1)
	c.Set("b", 2)

	c.Set("a", 10)

	c.Set("c", 3)

	if _, ok := c.Get("b"); ok {
		t.Error("Expected 'b' to be evicted")
	}
	if v, ok := c.Get("a"); !ok || v != 10 {
		t.Errorf("Expected 'a'=10, got ok=%v, v=%v", ok, v)
	}
}

func TestLRU_GetMissing(t *testing.T) {
	c := New[int](3)
	v, ok := c.Get("missing")
	if ok || v != 0 {
		t.Errorf("Expected (0, false), got (%v, %v)", v, ok)
	}
}

func TestLRU_Clear(t *testing.T) {
	c := New[int](3)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Clear()

	if _, ok := c.Get("a"); ok {
		t.Error("Expected cache to be empty after Clear")
	}
	if _, ok := c.Get("b"); ok {
		t.Error("Expected cache to be empty after Clear")
	}
}

func TestLRU_ConcurrentAccess(t *testing.T) {
	c := New[int](10)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i%10)
			c.Set(key, i)
			c.Get(key)
		}(i)
	}

	wg.Wait()
}
