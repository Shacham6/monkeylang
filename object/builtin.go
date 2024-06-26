package object

type BuiltinFunction func(...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (*Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (*Builtin) Inspect() string {
	return "builtin function"
}
