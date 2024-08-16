package queue

type Node[T comparable] struct {
	value T
	next  *Node[T]
	prev  *Node[T]
}

type Queue[T comparable] struct {
	head   *Node[T]
	tail   *Node[T]
	length int
}

func New[T comparable]() *Queue[T] {
	return &Queue[T]{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

func (q *Queue[T]) Len() int {
	return q.length
}

func (q *Queue[T]) Enqueue(value T) {
	node := &Node[T]{value: value}

	if q.head == nil {
		q.head = node
		q.tail = node
	} else {
		q.tail.next = node
		node.prev = q.tail
	}

	q.tail = node
	q.length++
}

func (q *Queue[T]) Dequeue() (val T, ok bool) {
	if q.head == nil {
		return val, false
	}

	node := q.head
	q.head = q.head.next
	q.length--

	if q.head == nil {
		q.tail = nil
	} else {
		q.head.prev = nil
	}

	return node.value, true
}

func (q *Queue[T]) Peek() (val T, ok bool) {
	if q.head == nil {
		return val, false
	}
	return q.head.value, true
}
