package general

import "errors"

// ErrStackEmpty is returned when attempting to access an element from an empty stack.
var ErrStackEmpty = errors.New("stack is empty")

// Stack implements the abstract data type (ADT) of a stack.
//
// Using a minimal set of operations (Push and Pop), it provides LIFO (last in, first out) access
// to the elements stored in it.
type Stack[T any] struct {
	slice []T
}

// IsEmpty checks if the stack is empty.
func (s Stack[T]) isEmpty() bool {
	return len(s.slice) == 0
}

// Top returns the top element of the stack without removing it.
//
// Returns [ErrStackEmpty] if the stack is empty.
func (s Stack[T]) top() (*T, error) {
	if s.isEmpty() {
		return nil, ErrStackEmpty
	}
	return &s.slice[len(s.slice)-1], nil
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.slice = append(s.slice, v)
}

// Pop removes and returns the top element of the stack.
//
// Returns [ErrStackEmpty] if the stack is empty.
func (s *Stack[T]) Pop() (*T, error) {
	top, err := s.top()
	if err != nil {
		return nil, err
	}
	s.slice = s.slice[:len(s.slice)-1]
	return top, nil
}
