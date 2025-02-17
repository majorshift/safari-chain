package types

import (
	"fmt"
	"reflect"
)

// List is a generic list implementation
// that stores elements of any type
type List[T any] struct {
	Data []T
}

func NewList[T any]() *List[T] {
	return &List[T]{
		Data: []T{},
	}
}

// Get finds element by index
func (l *List[T]) Get(index int) T {
	if index > len(l.Data)-1 {
		err := fmt.Sprintf("the given index (%d) is higher than the length (%d)", index, len(l.Data))
		panic(err)
	}
	return l.Data[index]
}

// Insert adds item to list
func (l *List[T]) Insert(v T) {
	l.Data = append(l.Data, v)
}

// Clear deletes all items in the list
func (l *List[T]) Clear() {
	l.Data = []T{}
}

// GetIndex will return the index of v. If v does not exist in the list
// -1 will be returned.
func (l *List[T]) GetIndex(v T) int {
	for i := 0; i < l.Len(); i++ {
		if reflect.DeepEqual(v, l.Data[i]) {
			return i
		}
	}
	return -1
}

// Remove deletes a item from the list
func (l *List[T]) Remove(v T) {
	index := l.GetIndex(v)
	if index == -1 {
		return
	}
	l.Pop(index)
}

// Pop deletes item matching index
func (l *List[T]) Pop(index int) {
	l.Data = append(l.Data[:index], l.Data[index+1:]...)
}

// Contains checks if an item v exists in the list
func (l *List[T]) Contains(v T) bool {
	for i := 0; i < len(l.Data); i++ {
		if reflect.DeepEqual(l.Data[i], v) {
			return true
		}
	}
	return false
}

// Last returns the last item in the list
func (l List[T]) Last() T {
	return l.Data[l.Len()-1]
}

// Len returns the size of the list
func (l *List[T]) Len() int {
	return len(l.Data)
}
