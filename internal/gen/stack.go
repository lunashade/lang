package gen

type Stack[T any] []*T

func NewStack[T any]() *Stack[T] {
	v := make(Stack[T], 0)
	return &v
}

func (s *Stack[T]) Push(x *T) {
	(*s) = append((*s), x)
}

func (s *Stack[T]) Pop() *T {
	v := (*s)[len(*s)-1]
	(*s)[len(*s)-1] = nil
	(*s) = (*s)[:len(*s)-1]
	return v
}

func (s *Stack[T]) Top() *T {
	return (*s)[len(*s)-1]
}
