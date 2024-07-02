package evaluator

import (
	"fmt"
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newWrongNumOfArgsError(len(args), 1)
			}

			switch arg := args[0].(type) {

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newWrongNumOfArgsError(len(args), 1)
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be an ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return &NULL
		},
	},

	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newWrongNumOfArgsError(len(args), 1)
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be an ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return &NULL
		},
	},

	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newWrongNumOfArgsError(len(args), 1)
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be an ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return &NULL
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newWrongNumOfArgsError(len(args), 2)
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("first argument to `push` must be %s, got %s", object.ARRAY_OBJ, args[0].Type())
			}

			arr := args[0].(*object.Array)
			return &object.Array{Elements: append(arr.Elements, args[1])}
		},
	},

	"puts": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newWrongNumOfArgsError(len(args), 1)
			}

			if args[0].Type() != object.STRING_OBJ {
				return newError("first argument to `puts` must be %s, got %s", object.STRING_OBJ, args[0].Type())
			}

			fmt.Println(args[0].Inspect())
			return &NULL
		},
	},
}

func newWrongNumOfArgsError(got int, want int64) *object.Error {
	return newError("wrong number of arguments. got = %d, want = %d", got, want)
}
