package linkedlist

import (
	"testing"
)

func TestIntLinkedList(t *testing.T) {
	list := New[int]()

	list.Append(5)
	list.Append(7)
	list.Append(9)

	test, ok := list.GetAt(2)
	if test != 9 || !ok {
		t.Errorf("Expected 9")
	}

	list.RemoveAt(1)

	list.RemoveAt(0)

	list.RemoveAt(0)

	if list.Len() != 0 {
		t.Errorf("Expected length of 0")
	}

	list.Prepend(5)
	list.Prepend(7)
	list.Prepend(9)

	test, ok = list.GetAt(0)
	if test != 9 || !ok {
		t.Errorf("Expected 9")
	}

	test, ok = list.GetAt(2)
	if test != 5 || !ok {
		t.Errorf("Expected 5")
	}

	list.Remove(9)

	if list.Len() != 2 {
		t.Errorf("Expected length of 2, got %d", list.Len())
	}

	test, ok = list.GetAt(0)
	if test != 7 || !ok {
		t.Errorf("Expected 7")
	}
}
