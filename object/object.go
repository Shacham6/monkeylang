package object

type Object interface {
	Type() ObjectType
	Inspect() string

	// Does not explicitly return the Hash object because some objects are not hashable,
	// and so we want to maintain the ability to return an error instead.
	HashKey() (HashKey, *Error)
}
