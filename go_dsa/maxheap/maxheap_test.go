package maxheap

import (
	"testing"
)

func TestIntMaxHeap(t *testing.T) {
	heap := New[int]()

	if heap.Len() != 0 {
		t.Errorf("Expected length to be 0, got %d", heap.Len())
	}

	heap.Insert(5)
	heap.Insert(3)
	heap.Insert(69)
	heap.Insert(420)
	heap.Insert(4)
	heap.Insert(1)
	heap.Insert(8)
	heap.Insert(7)

	if heap.Len() != 8 {
		t.Errorf("Expected length to be 8, got %d", heap.Len())
	}

	test, ok := heap.Delete()
	if test != 420 || !ok {
		t.Errorf("Expected 1")
	}

	test, ok = heap.Delete()
	if test != 69 || !ok {
		t.Errorf("Expected 3")
	}

	test, ok = heap.Delete()
	if test != 8 || !ok {
		t.Errorf("Expected 4")
	}

	test, ok = heap.Delete()
	if test != 7 || !ok {
		t.Errorf("Expected 5")
	}

	if heap.Len() != 4 {
		t.Errorf("Expected length to be 4, got %d", heap.Len())
	}

	test, ok = heap.Delete()
	if test != 5 || !ok {
		t.Errorf("Expected 7")
	}

	test, ok = heap.Delete()
	if test != 4 || !ok {
		t.Errorf("Expected 8")
	}

	test, ok = heap.Delete()
	if test != 3 || !ok {
		t.Errorf("Expected 69")
	}

	test, ok = heap.Delete()
	if test != 1 || !ok {
		t.Errorf("Expected 420")
	}

	if heap.Len() != 0 {
		t.Errorf("Expected length to be 0, got %d", heap.Len())
	}
}
