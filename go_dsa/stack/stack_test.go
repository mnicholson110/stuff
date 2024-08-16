package stack

import (
	"testing"
)

func TestIntStack(t *testing.T) {
	s := New[int]()

	s.Push(5)
	s.Push(7)
	s.Push(9)

	test, ok := s.Pop()
	if test != 9 || !ok {
		t.Error("Expected 9")
	}

	if s.Len() != 2 {
		t.Error("Expected length 2")
	}

	s.Push(11)

	test, ok = s.Pop()
	if test != 11 || !ok {
		t.Error("Expected 11")
	}

	test, ok = s.Pop()
	if test != 7 || !ok {
		t.Error("Expected 7")
	}

	test, ok = s.Peek()
	if test != 5 || !ok {
		t.Error("Expected 5")
	}

	test, ok = s.Pop()
	if test != 5 || !ok {
		t.Error("Expected 5")
	}

	if s.Len() != 0 {
		t.Error("Expected length 0")
	}

	s.Push(69)

	test, ok = s.Peek()
	if test != 69 || !ok {
		t.Error("Expected 69")
	}

	if s.Len() != 1 {
		t.Error("Expected length 1")
	}
}
