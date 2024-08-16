package stack

type Node[T comparable] struct {
	value T
	next  *Node[T]
	prev  *Node[T]
}

type Stack[T comparable] struct {
	head   *Node[T]
	tail   *Node[T]
	length int
}

func New[T comparable]() *Stack[T] {
	return &Stack[T]{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

func (s *Stack[T]) Len() int {
	return s.length
}

func (s *Stack[T]) Push(value T) {
	node := &Node[T]{value: value}
	if s.head == nil {
		s.head = node
		s.tail = node
	} else {
		s.head.prev = node
		node.next = s.head
		s.head = node
	}
	s.length++
}

func (s *Stack[T]) Pop() (val T, ok bool) {
	if s.head == nil {
		return val, false
	}

	node := s.head
	s.head = s.head.next
	s.length--

	if s.head == nil {
		s.tail = nil
	} else {
		s.head.prev = nil
	}

	return node.value, true
}

func (s *Stack[T]) Peek() (val T, ok bool) {
	if s.head == nil {
		return val, false
	}
	return s.head.value, true
}
