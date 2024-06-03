package object

type Object interface {
	Type() ObjectType
	Inspect() string
}
