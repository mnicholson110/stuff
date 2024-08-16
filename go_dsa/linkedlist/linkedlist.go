package linkedlist

type Node[T comparable] struct {
	value T
	next  *Node[T]
	prev  *Node[T]
}

type LinkedList[T comparable] struct {
	head   *Node[T]
	tail   *Node[T]
	length int
}

func New[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

func (l *LinkedList[T]) Len() int {
	return l.length
}

func (l *LinkedList[T]) Append(value T) {
	node := &Node[T]{value: value}

	if l.head == nil {
		l.head = node
	} else {
		l.tail.next = node
		node.prev = l.tail
	}

	l.tail = node
	l.length++
}

func (l *LinkedList[T]) Prepend(value T) {
	node := &Node[T]{value: value}

	if l.head == nil {
		l.tail = node
	} else {
		l.head.prev = node
		node.next = l.head
	}

	l.head = node
	l.length++
}

func (l *LinkedList[T]) AddAt(index int, value T) {
	if index < 0 || index > l.length {
		return
	}

	if index == l.length {
		l.Append(value)
		return
	}

	if index == 0 {
		l.Prepend(value)
		return
	}

	node := &Node[T]{value: value}
	current := l.head
	for i := 0; i < index; i++ {
		current = current.next
	}

	node.next = current
	node.prev = current.prev
	current.prev.next = node
	current.prev = node
	l.length++
}

func (l *LinkedList[T]) Remove(value T) {
	current := l.head
	for current != nil {
		if current.value == value {
			if current == l.head {
				l.head = l.head.next
				if l.head != nil {
					l.head.prev = nil
				}
			} else if current == l.tail {
				l.tail = l.tail.prev
				l.tail.next = nil
			} else {
				current.prev.next = current.next
				current.next.prev = current.prev
			}

			l.length--
			return
		}

		current = current.next
	}
}

func (l *LinkedList[T]) RemoveAt(index int) {
	if index < 0 || index >= l.length {
		return
	}

	if index == 0 {
		l.head = l.head.next
		if l.head != nil {
			l.head.prev = nil
		}
	} else if index == l.length-1 {
		l.tail = l.tail.prev
		l.tail.next = nil
	} else {
		current := l.head
		for i := 0; i < index; i++ {
			current = current.next
		}

		current.prev.next = current.next
		current.next.prev = current.prev
	}

	l.length--
}

func (l *LinkedList[T]) GetAt(index int) (val T, ok bool) {
	if index < 0 || index >= l.length {
		return val, false
	}

	current := l.head
	for i := 0; i < index; i++ {
		current = current.next
	}

	return current.value, true
}
