package object

import "fmt"

type BuiltinItem struct {
	Name    string
	Builtin *Builtin
}

func toBF(bf BuiltinFunction) *Builtin {
	return &Builtin{bf}
}

var Builtins = func() []BuiltinItem {
	return []BuiltinItem{
		{
			"len",
			toBF(func(args ...Object) Object {
				if len(args) != 1 {
					return newWrongNumOfArgsError(len(args), 1)
				}

				switch arg := args[0].(type) {

				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}

				case *String:
					return &Integer{Value: int64(len(arg.Value))}

				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			}),
		},
		{
			"first",
			toBF(func(args ...Object) Object {
				if len(args) != 1 {
					return newWrongNumOfArgsError(len(args), 1)
				}

				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `first` must be an ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*Array)
				if len(arr.Elements) > 0 {
					return arr.Elements[0]
				}

				return &CONST_NULL
			}),
		},
		{
			"last",
			toBF(func(args ...Object) Object {
				if len(args) != 1 {
					return newWrongNumOfArgsError(len(args), 1)
				}

				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `last` must be an ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*Array)
				length := len(arr.Elements)
				if length > 0 {
					return arr.Elements[length-1]
				}
				return &Null{}
			}),
		},
		{
			"rest",
			toBF(func(args ...Object) Object {
				if len(args) != 1 {
					return newWrongNumOfArgsError(len(args), 1)
				}

				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `rest` must be an ARRAY, got %s", args[0].Type())
				}

				arr := args[0].(*Array)
				length := len(arr.Elements)
				if length > 0 {
					newElements := make([]Object, length-1)
					copy(newElements, arr.Elements[1:length])
					return &Array{Elements: newElements}
				}

				return &CONST_NULL
			}),
		},
		{
			"push",
			toBF(func(args ...Object) Object {
				if len(args) != 2 {
					return newWrongNumOfArgsError(len(args), 2)
				}

				if args[0].Type() != ARRAY_OBJ {
					return newError("first argument to `push` must be %s, got %s", ARRAY_OBJ, args[0].Type())
				}

				arr := args[0].(*Array)
				return &Array{Elements: append(arr.Elements, args[1])}
			}),
		},
		{
			"puts",
			toBF(func(args ...Object) Object {
				if len(args) != 1 {
					return newWrongNumOfArgsError(len(args), 1)
				}

				if args[0].Type() != STRING_OBJ {
					return newError("first argument to `puts` must be %s, got %s", STRING_OBJ, args[0].Type())
				}

				fmt.Println(args[0].Inspect())
				return &CONST_NULL
			}),
		},
		{
			"sprintf",
			toBF(func(args ...Object) Object {
				if len(args) < 1 {
					return newError("sprintf function requires at least a single argument")
				}

				fmtArg := args[0]
				if fmtArg.Type() != STRING_OBJ {
					return newError("first argument to `sprintf` must be %s, got = %s", STRING_OBJ, fmtArg.Type())
				}

				formattedArgs := make([]any, len(args)-1)
				for i := 0; i < len(args)-1; i++ {
					formattedArgs[i] = args[i+1].Inspect()
				}

				return &String{fmt.Sprintf(fmtArg.Inspect(), formattedArgs...)}
			}),
		},
	}
}()

func newError(format string, args ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, args...)}
}

func newWrongNumOfArgsError(got int, want int64) *Error {
	return newError("wrong number of arguments. got = %d, want = %d", got, want)
}
