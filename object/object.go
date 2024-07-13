package object

type Object interface {
	hashKey
	deval

	Type() ObjectType
	Inspect() string
}
