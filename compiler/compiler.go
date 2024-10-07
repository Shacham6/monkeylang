package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object

	lastInstruction EmittedInstruction
	prevInstruction EmittedInstruction
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},

		lastInstruction: EmittedInstruction{},
		prevInstruction: EmittedInstruction{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
		return nil

	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		c.emit(code.OpPop)
		return nil

	case *ast.InfixExpression:
		// this operator is treated as equivalent to others *just with flipped operands*.
		// in essence: `1 < 2` is translated into `2 > 1`.
		if node.Operator == "<" {
			if err := c.Compile(node.Right); err != nil {
				return err
			}

			if err := c.Compile(node.Left); err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)
			return nil
		}

		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Right); err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("infix operator %s not supported", node.Operator)
		}
		return nil

	case *ast.Boolean:
		if node.Value() {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
		return nil

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
		return nil

	case *ast.PrefixExpression:
		if err := c.Compile(node.Right); err != nil {
			return err
		}

		switch node.Operator {
		case "-":
			c.emit(code.OpMinus)
		case "!":
			c.emit(code.OpBang)
		default:
			panic(fmt.Sprintf("prefix operator %s is not supported", node.Operator))
		}
		return nil

	case *ast.IfExpression:
		if err := c.Compile(node.Condition()); err != nil {
			return err
		}

		// Emit the opcode with a bogus offset
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)
		if err := c.Compile(node.Consequence()); err != nil {
			return err
		}

		if c.lastInstruction.Opcode == code.OpPop {
			c.removeLastInstruction()
		}

		alternative, hasAlternative := node.Alternative()

		if !hasAlternative {
			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)
		} else {
			// Emit an 'OpJump' with a bogus value
			jumpPos := c.emit(code.OpJump, 9999)

			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

			if err := c.Compile(alternative); err != nil {
				return err
			}

			if c.lastInstruction.Opcode == code.OpPop {
				c.removeLastInstruction()
			}
			c.changeOperand(jumpPos, len(c.instructions))
		}

		return nil

	case *ast.BlockStatement:
		for _, s := range node.Statements() {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
		return nil

	default:
		panic(fmt.Sprintf("don't support node of type %T", node))
	}
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		c.instructions,
		c.constants,
	}
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	c.updateLastInstruction(op, pos)
	return pos
}

func (c *Compiler) updateLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{
		Opcode:   op,
		Position: pos,
	}
	c.prevInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) removeLastInstruction() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.prevInstruction
}

// replaceInstruction replaces an instruction, defined at pos, with the newInstruction.
//
// Please be careful using this. Unless the instruction at pos and the new instruction
// are of the same size, you are risking an overflow. Tread carefully.
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

// changeOperand changes the operand of the instruction at opPos.
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}
