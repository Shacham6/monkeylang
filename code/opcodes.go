package code

import "fmt"

const (
	OpConstant Opcode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPop
	OpNull
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
	OpGetGlobal
	OpSetGlobal
	OpArray
	OpHash
	OpIndex
	OpCall
	OpReturnValue
	OpReturn
	OpGetLocal
	OpSetLocal
	// This is an element of the book I don't think I can entirely agree with in terms of language design, right?
	// I mean it can't be a good thing to have an entirely separate way to handle the modules of just the standard
	// library, after all a shitload of functions are called all the time.
	// ...Or is it more logical that by enhancing these functions, everything's enhanced, because everything must
	// end up calling these functions, but even so if I _can_ do this it's because I statically know their identity
	// and their memory location and can thus optimize for such, can't I do this for all functions?
	// I don't know man.
	OpGetBuiltin

	OpClosure
	OpGetFree

	OpCurrentClosure
)

var definitions = map[Opcode]*Definition{
	OpConstant:       {"OpConstant", []int{2}},
	OpAdd:            {"OpAdd", []int{}},
	OpSub:            {"OpSub", []int{}},
	OpMul:            {"OpMul", []int{}},
	OpDiv:            {"OpDiv", []int{}},
	OpPop:            {"OpPop", []int{}},
	OpNull:           {"OpNull", []int{}},
	OpTrue:           {"OpTrue", []int{}},
	OpFalse:          {"OpFalse", []int{}},
	OpEqual:          {"OpEqual", []int{}},
	OpNotEqual:       {"OpNotEqual", []int{}},
	OpGreaterThan:    {"OpGreaterThan", []int{}},
	OpMinus:          {"OpMinus", []int{}},
	OpBang:           {"OpBang", []int{}},
	OpJumpNotTruthy:  {"OpJumpNotTruthy", []int{2}},
	OpJump:           {"OpJump", []int{2}},
	OpGetGlobal:      {"OpGetGlobal", []int{2}},
	OpSetGlobal:      {"OpSetGlobal", []int{2}},
	OpArray:          {"OpArray", []int{2}}, // operand here is 2 bytes wide, which gives us 65535 possible number of elements
	OpHash:           {"OpHash", []int{2}},  // operand here is 2 bytes wide, which gives us 65535 possible number of elements
	OpIndex:          {"OpIndex", []int{}},
	OpCall:           {"OpCall", []int{1}},
	OpReturnValue:    {"OpReturnValue", []int{}},
	OpReturn:         {"OpReturn", []int{}},
	OpGetLocal:       {"OpGetLocal", []int{1}},
	OpSetLocal:       {"OpSetLocal", []int{1}},
	OpGetBuiltin:     {"OpGetBuiltin", []int{1}},
	OpClosure:        {"OpClosure", []int{2, 1}},
	OpGetFree:        {"OpGetFree", []int{1}},
	OpCurrentClosure: {"OpCurrentClosure", []int{}},
}

type Definition struct {
	Name string

	// OperandWidths is the widths of the different operands.
	//
	// Each element represents a different operand, and the value is its size
	// in bytes.
	// An example value of "[]int{2}" means that the Definition has a single
	// operand, sized at 16 bytes (or 2x8).
	OperandWidths []int
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}
