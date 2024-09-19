package compiler

import (
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
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

	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.InfixExpression:
		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Right); err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		default:
			return fmt.Errorf("infix operator %s not supported", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	default:
		panic(fmt.Sprintf("don't support node of type %T", node))
	}
	return nil
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
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}
