package maxheap

import (
	"golang.org/x/exp/constraints"
)

type MaxHeap[T constraints.Ordered] struct {
	heap   []T
	length int
}

func New[T constraints.Ordered]() *MaxHeap[T] {
	return &MaxHeap[T]{
		heap:   make([]T, 0),
		length: 0,
	}
}

func (h *MaxHeap[T]) Len() int {
	return h.length
}

func (h *MaxHeap[T]) Insert(value T) {
	h.heap = append(h.heap, value)
	h.heapifyUp(h.length)
	h.length++
}

func (h *MaxHeap[T]) Delete() (val T, ok bool) {
	if h.length == 0 {
		return val, false
	}

	out := h.heap[0]

	if h.length == 1 {
		h.length--
		h.heap = h.heap[:0]
		return out, true
	}

	h.length--
	h.heap[0] = h.heap[h.length]
	h.heapifyDown(0)

	return out, true
}

func (h *MaxHeap[T]) heapifyUp(idx int) {
	if idx == 0 {
		return
	}

	parentIdx := h.parentIdx(idx)

	if h.heap[parentIdx] < h.heap[idx] {
		h.heap[parentIdx], h.heap[idx] = h.heap[idx], h.heap[parentIdx]
		h.heapifyUp(parentIdx)
	}
}

func (h *MaxHeap[T]) heapifyDown(idx int) {
	leftChildIdx := h.leftChildIdx(idx)
	rightChildIdx := h.rightChildIdx(idx)

	if idx >= h.length || leftChildIdx >= h.length {
		return
	}

	if h.heap[leftChildIdx] < h.heap[rightChildIdx] && h.heap[idx] < h.heap[rightChildIdx] {
		h.heap[idx], h.heap[rightChildIdx] = h.heap[rightChildIdx], h.heap[idx]
		h.heapifyDown(rightChildIdx)
	} else if h.heap[idx] < h.heap[leftChildIdx] {
		h.heap[idx], h.heap[leftChildIdx] = h.heap[leftChildIdx], h.heap[idx]
		h.heapifyDown(leftChildIdx)
	}
}

func (h *MaxHeap[T]) parentIdx(idx int) int {
	return (idx - 1) / 2
}

func (h *MaxHeap[T]) leftChildIdx(idx int) int {
	return 2*idx + 1
}

func (h *MaxHeap[T]) rightChildIdx(idx int) int {
	return 2*idx + 2
}
