package q

import "testing"

func TestLRUCache(t *testing.T) {
	c := NewLRUCache(2)

	c.Put(1, 10)
	c.Put(2, 20)
	if val := c.Get(1); val != 10 {
		t.Fatalf("got %d want %d", val, 10)
	}

	c.Put(3, 30)
	if val := c.Get(2); val != -1 {
		t.Fatalf("got %d want %d", val, -1)
	}

	c.Put(4, 40)
	if val := c.Get(1); val != -1 {
		t.Fatalf("got %d want %d", val, -1)
	}
	if val := c.Get(3); val != 30 {
		t.Fatalf("got %d want %d", val, 30)
	}
}
