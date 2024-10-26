package unsafestack

// UnsafeSizedStack is a stack data structure that prioritizes
// speed over correctness in a couple of ways. It it made to be
// generic.
//
// The stack will preallocate storage ahead of time, and operate
// a constantly moving pointer. For performance's sake, bounds
// will not be checked.
//
// Also, popped items are not actually deleted. Instead, what happens
// is that we mark the area as writeable.
type UnsafeSizedStack[T any] struct {
	items []T

	// topIndex is the position of the top of the stack.
	//
	// This position can be used to signify the safest reading position.
	// This position+1 is the writing position.
	topIndex int
}

func Make[T any](size int) UnsafeSizedStack[T] {
	return UnsafeSizedStack[T]{
		items:    make([]T, size),
		topIndex: -1,
	}
}

func (u *UnsafeSizedStack[T]) Push(item T) {
	u.items[u.topIndex+1] = item
	u.topIndex++
}

func (u *UnsafeSizedStack[T]) Size() int {
	// doing the `+1` because an empty stack means that the reading position
	// sits at `topIndex == -1`. See the initialization at the `Make` func.
	return u.topIndex + 1
}

func (u *UnsafeSizedStack[T]) Pop() T {
	// instead of saving the item, we can just return the item from "out of bounds".
	// this is a niche advantage of not actually deleted the data.
	u.topIndex--
	return u.items[u.topIndex+1]
}
