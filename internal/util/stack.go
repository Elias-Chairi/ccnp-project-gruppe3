package util

// Stack implements the abstract data type of a stack.
// First In Last Out (FILO) principle.
type Stack[T any] struct{
	slice []T
}

// Top returns the top element of the stack without removing it.
func (s Stack[T]) Top() T {
	return s.slice[len(s.slice)-1]
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.slice = append(s.slice, v)
}

// Pop removes and returns the top element of the stack.
func (s *Stack[T]) Pop() T{
	top := s.Top()
	s.slice = s.slice[:len(s.slice)-1]
	return top
}

























