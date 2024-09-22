package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 2048

type VM struct {
	// parts of the bytescode
	constants    []object.Object
	instructions code.Instructions

	// mutating runtime things
	stack []object.Object
	sp    int // Always points to the next value. Top of stack is stack[sp-1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		bytecode.Constants,
		bytecode.Instructions,

		make([]object.Object, StackSize),
		0,
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
	for ip := 0; ip < len(vm.instructions); ip++ {
		// we're on the *hot* path, this is the actual execution of the vm, thus
		// we're not using `code.Lookup` since it'll slow things down for us.
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			if err := vm.push(vm.constants[constIndex]); err != nil {
				return err
			}

		case code.OpAdd, code.OpSub, code.OpDiv, code.OpMul:
			if err := vm.executeBinaryOperation(op); err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		default:
			rawCode := vm.instructions[ip]
			definition, err := code.Lookup(rawCode)
			if err != nil {
				// TODO: think if this flow makes sense; if does add test.
				return fmt.Errorf("encountered an unknown opcode: %q", rawCode)
			}

			// TODO: I'm not sure if this really makes sense to test, consider
			// changing into a panic maybe?
			return fmt.Errorf("opcode %s not yet supported", definition.Name)
		}
	}

	return nil
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

	return fmt.Errorf(
		"unsupported types for binary operations: %s %s",
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

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp-- // simply decreasing the pointer, this will allow this location in memory to be overwritten. no need to explicitly "drop" the memory.
	return o
}
