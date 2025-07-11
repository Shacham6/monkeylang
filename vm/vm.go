package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
	unsafestack "monkey/unsafe_stack"
)

const (
	StackSize   = 2048
	GlobalsSize = 65536
	MaxFrames   = 1024
)

func InitGlobalsArray() []object.Object {
	return make([]object.Object, GlobalsSize)
}

func InitStackArray() []object.Object {
	return make([]object.Object, StackSize)
}

func InitFramesArray() []*Frame {
	return make([]*Frame, MaxFrames)
}

var (
	constTrue  = &object.Boolean{Value: true}
	constFalse = &object.Boolean{Value: false}
	constNull  = &object.Null{}
)

type VM struct {
	// parts of the bytescode
	constants []object.Object

	// mutating runtime things
	stack []object.Object
	sp    int // "stack pointer". Always points to the next value. Top of stack is stack[sp-1]

	globals []object.Object

	frameStack unsafestack.UnsafeSizedStack[*Frame]
}

func New(bytecode *compiler.Bytecode) *VM {
	return NewWithGlobalState(
		bytecode,
		InitGlobalsArray(),
	)
}

func NewWithGlobalState(bytecode *compiler.Bytecode, globals []object.Object) *VM {
	sp := 0

	framesStack := unsafestack.Make[*Frame](MaxFrames)

	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	framesStack.Push(mainFrame)

	return &VM{
		bytecode.Constants,
		InitStackArray(),
		sp,
		globals,
		framesStack,
	}
}

// StackTop returns the top object of the stack without popping.
//
// if the stack is empty, nil will be returned.
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

// LastPoppedStackElem returns the "overflown" value from the stack,
// suggesting that the value was popped itself.
//
// This operation is not safe, and if nothing was popped and this
// method called - the program will crash.
//
// The reason this works is that "popping" the stack does not explicitly
// delete the data. Instead, the pointer marking length of written data
// is decremented, signaling that this section of the stack is now writeable.
func (vm *VM) LastPoppedStackElem() object.Object {
	if vm.sp >= len(vm.stack) {
		panic(
			"Attempting to view last popped stack elem before anything was ever popped, this state should be impossible!",
		)
	}
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	var ip int
	var ins code.Instructions
	var op code.Opcode

	toErr := func(err error) error {
		return &VmRunError{err, ins, vm.stack, vm.globals, vm.sp}
	}

	for vm.frameStack.Current().ip < len(vm.frameStack.Current().Instructions())-1 {
		// we're on the *hot* path, this is the actual execution of the vm, thus
		// we're not using `code.Lookup` since it'll slow things down for us.

		vm.frameStack.Current().ip++
		ip = vm.frameStack.Current().ip
		ins = vm.frameStack.Current().Instructions()
		op = code.Opcode(ins[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.frameStack.Current().ip += 2
			if err := vm.push(vm.constants[constIndex]); err != nil {
				return toErr(err)
			}

		case code.OpJump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.frameStack.Current().ip = pos - 1

		case code.OpJumpNotTruthy:
			obj := vm.pop()
			if objectBoolToNativeBool(obj) {
				// We add by the width (in bytes) of the operands.
				vm.frameStack.Current().ip += 2

				continue
			}

			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.frameStack.Current().ip = pos - 1

		case code.OpAdd, code.OpSub, code.OpDiv, code.OpMul:
			if err := vm.executeBinaryOperation(op); err != nil {
				return toErr(err)
			}

		case code.OpPop:
			vm.pop()

		case code.OpTrue:
			if err := vm.push(constTrue); err != nil {
				return toErr(err)
			}

		case code.OpFalse:
			if err := vm.push(constFalse); err != nil {
				return toErr(err)
			}

		case code.OpNull:
			if err := vm.push(constNull); err != nil {
				return toErr(err)
			}

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			if err := vm.executeComparison(op); err != nil {
				return toErr(err)
			}

		case code.OpBang:
			if err := vm.executeBangOperator(); err != nil {
				return toErr(err)
			}

		case code.OpMinus:
			if err := vm.executeMinusOperator(); err != nil {
				return toErr(err)
			}

		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.frameStack.Current().ip += 2
			vm.globals[globalIndex] = vm.pop()

		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.frameStack.Current().ip += 2
			if err := vm.push(vm.globals[globalIndex]); err != nil {
				return toErr(err)
			}

		case code.OpGetBuiltin:
			builtinIndex := code.ReadUint8(ins[ip+1:])
			vm.frameStack.Current().ip += 1
			if err := vm.push(object.Builtins[builtinIndex].Builtin); err != nil {
				return toErr(err)
			}

		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.frameStack.Current().ip += 1

			frame := vm.frameStack.Current()
			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()

		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.frameStack.Current().ip += 1

			currentFrame := vm.frameStack.Current()
			local := vm.stack[currentFrame.basePointer+int(localIndex)]
			if err := vm.push(local); err != nil {
				return toErr(err)
			}

		case code.OpGetFree:
			index := code.ReadUint8(ins[ip+1:])
			vm.frameStack.Current().ip += 1
			currentClosure := vm.frameStack.Current().cl
			if err := vm.push(currentClosure.Free[index]); err != nil {
				return toErr(err)
			}

		case code.OpCurrentClosure:
			currentClosure := vm.frameStack.Current().cl
			if err := vm.push(currentClosure); err != nil {
				return toErr(err)
			}

		case code.OpArray:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.frameStack.Current().ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			if err := vm.push(array); err != nil {
				return toErr(err)
			}

		case code.OpHash:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.frameStack.Current().ip += 2

			hashmap := map[object.HashKey]object.HashPair{}
			for index := vm.sp - numElements; index < vm.sp; index += 2 {
				key := vm.stack[index]

				// postpone user error handling for a moment
				hashkey, err := key.HashKey()
				if err != nil {
					return toErr(fmt.Errorf("type is unusable as a hash key: %s", key.Type()))
				}

				value := vm.stack[index+1]

				hashmap[hashkey] = object.HashPair{
					Key:   key,
					Value: value,
				}
			}

			vm.sp = vm.sp - numElements

			if err := vm.push(&object.Hash{Pairs: hashmap}); err != nil {
				return toErr(err)
			}

		case code.OpIndex:
			index := vm.pop()
			collection := vm.pop()

			switch collection := collection.(type) {
			case *object.Hash:
				if err := vm.executeHashIndexOperator(collection, index); err != nil {
					return toErr(err)
				}

			case *object.Array:
				if err := vm.executeArrayIndexOperator(collection, index); err != nil {
					return toErr(err)
				}
			}

		case code.OpClosure:
			constIndex := code.ReadUint16(ins[ip+1:])
			numFree := code.ReadUint8(ins[ip+3:])
			vm.frameStack.Current().ip += 3

			err := vm.pushClosure(int(constIndex), int(numFree))
			if err != nil {
				return toErr(err)
			}

		case code.OpCall:
			numOfArgs := code.ReadUint8(ins[ip+1:])
			vm.frameStack.Current().ip += 1

			calleeT := vm.stack[vm.sp-1-int(numOfArgs)]

			iNumOfArgs := int(numOfArgs)

			// The fact that we're matching called arguments at runtime is weird for me,
			// given this info is available at compile time... But no matter for now.
			// Following along.
			switch callee := calleeT.(type) {
			case *object.Closure:
				if callee.Fn.NumParameters != iNumOfArgs {
					return toErr(fmt.Errorf(
						"wrong number of arguments: want = %d, got = %d",
						callee.Fn.NumParameters,
						iNumOfArgs,
					))
				}

				frame := NewFrame(callee, vm.sp-iNumOfArgs)
				vm.frameStack.Push(frame)
				vm.sp = frame.basePointer + callee.Fn.NumLocals

			case *object.Builtin:
				args := vm.stack[vm.sp-iNumOfArgs : vm.sp]
				result := callee.Fn(args...)
				vm.sp = vm.sp - iNumOfArgs - 1
				if result == nil {
					vm.push(constNull)
				} else {
					vm.push(result)
				}

			default:
				return toErr(fmt.Errorf(
					"calling a non-function: (%s) %s",
					callee.Type(),
					callee.Inspect(),
				))
			}

		case code.OpReturnValue:
			returnValue := vm.pop()

			// two pops - one for the function frame, and one for the CALL that
			// put us into the function to begin with.
			frame := vm.frameStack.Pop()
			vm.sp = frame.basePointer - 1

			if err := vm.push(returnValue); err != nil {
				return toErr(err)
			}

		case code.OpReturn:
			frame := vm.frameStack.Pop()
			vm.sp = frame.basePointer - 1
			if err := vm.push(constNull); err != nil {
				return toErr(err)
			}

		default:
			rawCode := ins[ip]
			definition, err := code.Lookup(rawCode)
			if err != nil {
				// TODO: think if this flow makes sense; if does add test.
				return toErr(fmt.Errorf("encountered an unknown opcode: %q", rawCode))
			}

			// TODO: I'm not sure if this really makes sense to test, consider
			// changing into a panic maybe?
			return toErr(fmt.Errorf("opcode %s not yet supported", definition.Name))
		}
	}

	return nil
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case constTrue:
		return vm.push(constFalse)
	case constFalse, constNull:
		return vm.push(constTrue)
	default:
		return vm.push(constFalse)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()
	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("prefix operator '-' not supported for type '%s'", operand.Type())
	}

	value := operand.(*object.Integer).Value

	return vm.push(&object.Integer{Value: -value})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToObjectBool(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBoolToObjectBool(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (vm *VM) buildArray(startIndex int, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	// I guess this means that the values must be arranged in memory.
	// ...And from what I've seen that's indeed the case? Yeah I think
	// so...
	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (vm *VM) executeHashIndexOperator(hash *object.Hash, index object.Object) error {
	// Hashing and stuff's reserved to the hashmap type.
	hashKey, err := index.HashKey()
	if err != nil {
		return fmt.Errorf("object %s is not hashable: %s", index.Type(), err)
	}

	if err := vm.push(hash.Pairs[hashKey].Value); err != nil {
		return err
	}

	return nil
}

func (vm *VM) executeArrayIndexOperator(array *object.Array, index object.Object) error {
	indexValue := index.(*object.Integer).Value
	return vm.push(array.Elements[indexValue])
}

func (vm *VM) executeIntegerComparison(
	op code.Opcode,
	left object.Object, right object.Object,
) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToObjectBool(leftValue == rightValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToObjectBool(leftValue != rightValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToObjectBool(leftValue > rightValue))
	default:
		panic(fmt.Sprintf("unexpected code.Opcode: %#v", op))
	}
	// if complaints about missing return, means that a branch is missing a return clause.
}

func nativeBoolToObjectBool(b bool) object.Object {
	if b {
		return constTrue
	}
	return constFalse
}

func objectBoolToNativeBool(o object.Object) bool {
	switch o.Type() {
	case object.BOOLEAN_OBJ:
		return o.(*object.Boolean).Value

	case object.INTEGER_OBJ:
		value := o.(*object.Integer)
		return value.Value != 0

	case object.NULL_OBJ:
		return false

	default:
		return true
	}
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	// @Shacham:
	// we're explicitly expecting only numbers here, will fail on floats.
	// on the other hand i'm pretty certain we don't support floats at all
	// as of now... so yeah.

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}

	if leftType == object.STRING_OBJ && rightType == object.STRING_OBJ {
		return vm.executeBinaryStringOperation(op, left, right)
	}

	def, err := code.Lookup(byte(op))
	if err != nil {
		panic(fmt.Sprintf("failed finding definition of the opcode: %s", err))
	}

	return fmt.Errorf(
		"unsupported types for binary (%s) operations: %s %s",
		def.Name,
		leftType,
		rightType,
	)
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left object.Object, right object.Object) error {
	leftInteger, ok := left.(*object.Integer)
	if !ok {
		return fmt.Errorf("%s an invalid %s", left.Inspect(), left.Type())
	}

	rightInteger, ok := right.(*object.Integer)
	if !ok {
		return fmt.Errorf("%s an invalid %s", right.Inspect(), right.Type())
	}

	leftValue, rightValue := leftInteger.Value, rightInteger.Value

	var result int64
	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unsupported integer operator: %d", op)
	}

	integerResult := object.Integer{Value: result}
	if err := vm.push(&integerResult); err != nil {
		return err
	}

	return nil
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left object.Object, right object.Object) error {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	switch op {
	case code.OpAdd:
		return vm.push(&object.String{Value: fmt.Sprintf("%s%s", leftValue, rightValue)})

	default:
		def, err := code.Lookup(byte(op))
		if err != nil {
			panic(fmt.Sprintf("opcode %q is not recognized: %s", byte(op), err))
		}
		return fmt.Errorf("binary operator %s not supported for strings", def.Name)
	}
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

func (vm *VM) pushClosure(constIndex int, numFree int) error {
	constant := vm.constants[constIndex]
	function, ok := constant.(*object.CompiledFunction)
	if !ok {
		return fmt.Errorf("invalid constant: expected COMPILED_FUNCTION but got %s", constant.Type())
	}

	free := make([]object.Object, numFree)
	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}
	// We move the pointer back so that the VM will move to execute the OpGet*
	// codes manually.
	vm.sp = vm.sp - numFree

	closure := &object.Closure{Fn: function, Free: free}
	return vm.push(closure)
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp-- // simply decreasing the pointer, this will allow this location in memory to be overwritten. No need to explicitly "drop" the memory.
	return o
}
