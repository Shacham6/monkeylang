package optional

type Optional[T any] struct {
	data *T
}

func Make[T any](content *T) Optional[T] {
	return Optional[T]{content}
}

func (o *Optional[T]) Ok() bool {
	return o.data != nil
}

func (o *Optional[T]) Data() *T {
	return o.data
}
