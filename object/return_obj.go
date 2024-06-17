package object

type ReturnValue struct {
	Value Object
}

func (rt *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rt *ReturnValue) Inspect() string {
	return rt.Value.Inspect()
}
