package compiler

import (
	"cmp"
	"fmt"
	"monkey/ast"
	"monkey/code"
	"monkey/object"
	"slices"
)

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func ZeroEmittedInstruction() EmittedInstruction {
	return EmittedInstruction{} //nolint:exhaustruct
}

type Compiler struct {
	constants []object.Object

	symbolTable *SymbolTable

	scopes     []CompilationScope
	scopeIndex int
}

func New() *Compiler {
	constants := []object.Object{}
	return NewWithState(NewSymbolTable(), constants)
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	mainScope := NewCompilationScope()

	return &Compiler{
		constants: constants,

		symbolTable: s,

		scopes:     []CompilationScope{mainScope},
		scopeIndex: 0,
	}
}

func (c *Compiler) scope() *CompilationScope {
	return &c.scopes[c.scopeIndex]
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

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))
		return nil

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			if err := c.Compile(el); err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(node.Elements))
		return nil

	case *ast.HashLiteral:
		// Go makes no guarantees regarding the order of the keys in maps,
		// we need to sort the keys in order to emit consistent bytecode.
		pairs := node.Pairs()
		keys := getKeys(pairs)
		slices.SortFunc(keys, func(a, b ast.Expression) int {
			return cmp.Compare(a.String(), b.String())
		})

		for _, key := range keys {
			value := pairs[key]

			if err := c.Compile(key); err != nil {
				return err
			}

			if err := c.Compile(value); err != nil {
				return err
			}
		}

		c.emit(code.OpHash, len(node.Pairs())*2)
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
		jumpNotTruthyInstuctionPos := c.emit(code.OpJumpNotTruthy, 9999)
		if err := c.Compile(node.Consequence()); err != nil {
			return err
		}

		if c.scope().LastInstructionIs(code.OpPop) {
			c.scope().RemoveLastInstruction()
		}

		// Emit the opcode with a bogus offset
		jumpInstructionPos := c.emit(code.OpJump, 9999)

		c.scope().ChangeOperand(
			jumpNotTruthyInstuctionPos,
			len(c.scope().Instructions),
		)

		alt, hasAlt := node.Alternative()
		if !hasAlt {
			c.emit(code.OpNull)
		} else {
			if err := c.Compile(alt); err != nil {
				return err
			}

			if c.scope().LastInstructionIs(code.OpPop) {
				c.scope().RemoveLastInstruction()
			}
		}
		c.scope().ChangeOperand(
			jumpInstructionPos,
			len(c.scope().Instructions),
		)

		return nil

	case *ast.BlockStatement:
		for _, s := range node.Statements() {
			if err := c.Compile(s); err != nil {
				return err
			}
		}
		return nil

	case *ast.Identifier:
		if node.Value == "null" {
			c.emit(code.OpNull)
			return nil
		}

		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable: %s", node.Value)
		}

		c.emit(code.OpGetGlobal, symbol.Index)
		return nil

	case *ast.LetStatement:
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, symbol.Index)
		return nil

	case *ast.IndexExpression:
		if err := c.Compile(node.Left()); err != nil {
			return err
		}

		if err := c.Compile(node.Index()); err != nil {
			return err
		}

		c.emit(code.OpIndex)
		return nil

	case *ast.FunctionLiteral:
		// ========== ENTER FUNCTION SCOPE ==========

		c.enterScope()

		if err := c.Compile(node.Body()); err != nil {
			return err
		}

		// This is that all "the last expression in the function is the return value"
		if c.scope().LastInstructionIs(code.OpPop) {
			c.scope().ReplaceInstruction(
				c.scope().LastInstruction.Position,
				code.Make(code.OpReturnValue),
			)
		}

		if !c.scope().LastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		instructions := c.leaveScope()

		// ========== LEAVING FUNCTION SCOPE ==========

		compiledFn := &object.CompiledFunction{Instructions: instructions}

		c.emit(code.OpConstant, c.addConstant(compiledFn))

		return nil

	case *ast.ReturnStatement:
		if err := c.Compile(node.ReturnValue); err != nil {
			return err
		}
		c.emit(code.OpReturnValue)
		return nil

	default:
		panic(fmt.Sprintf("don't support node of type %T", node))
	}
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		c.scope().Instructions,
		c.constants,
	}
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	return c.scope().Emit(op, operands...)
}

func (c *Compiler) enterScope() {
	scope := NewCompilationScope()
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.scope().Instructions

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	return instructions
}
