package queue

import (
	"testing"
)

func TestIntQueue(t *testing.T) {
	s := New[int]()

	s.Enqueue(5)
	s.Enqueue(7)
	s.Enqueue(9)

	test, ok := s.Dequeue()
	if test != 5 || !ok {
		t.Error("Expected 5")
	}

	if s.Len() != 2 {
		t.Error("Expected length 2")
	}

	s.Enqueue(11)

	test, ok = s.Dequeue()
	if test != 7 || !ok {
		t.Error("Expected 7")
	}

	test, ok = s.Dequeue()
	if test != 9 || !ok {
		t.Error("Expected 9")
	}

	test, ok = s.Peek()
	if test != 11 || !ok {
		t.Error("Expected 11")
	}

	test, ok = s.Dequeue()
	if test != 11 || !ok {
		t.Error("Expected 11")
	}

	if s.Len() != 0 {
		t.Error("Expected length 0")
	}

	s.Enqueue(69)

	test, ok = s.Peek()
	if test != 69 || !ok {
		t.Error("Expected 69")
	}

	if s.Len() != 1 {
		t.Error("Expected length 1")
	}
}
