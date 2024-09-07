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

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++
	return nil
}
