package lrucache

import (
	"testing"
)

func TestIntLRUCache(t *testing.T) {
	c := New[string, int](3)

	if _, r := c.Get("foo"); r {
		t.Error("Expected nil")
	}

	c.Update("foo", 69)

	v, r := c.Get("foo")
	if v != 69 || !r {
		t.Error("Expected 69")
	}

	c.Update("bar", 420)

	v, r = c.Get("bar")
	if v != 420 || !r {
		t.Error("Expected 420")
	}

	c.Update("baz", 1337)

	v, r = c.Get("baz")
	if v != 1337 || !r {
		t.Error("Expected 1337")
	}

	c.Update("ball", 69420)

	v, r = c.Get("ball")
	if v != 69420 || !r {
		t.Error("Expected 69420")
	}

	if _, r := c.Get("foo"); r {
		t.Error("Expected nil")
	}

	v, r = c.Get("bar")
	if v != 420 || !r {
		t.Error("Expected 420")
	}

	c.Update("foo", 69)

	if v != 420 || !r {
		t.Error("Expected 420")
	}

	v, r = c.Get("foo")
	if v != 69 || !r {
		t.Error("Expected 69")
	}

	if _, r := c.Get("baz"); r {
		t.Error("Expected nil")
	}
}
