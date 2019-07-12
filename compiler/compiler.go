package compiler

import (
	"azula/ast"
	"azula/code"
	"azula/object"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"fmt"
)

var typeMap = map[object.Type]string{
	object.IntegerObj: "int",
	object.BooleanObj: "bool",
	object.StringObj:  "string",
	object.ArrayObj:   "array",
	object.ErrorObj:   "error",
}

var azulaTypesToLLVM = map[string]types.Type{
	"int":    types.I32,
	"string": types.NewArray(1024, types.I8),
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type Compiler struct {
	Module   *ir.Module
	CurBlock *ir.Block

	Functions map[string]*ir.Func
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		Module:    ir.NewModule(),
		CurBlock:  nil,
		Functions: make(map[string]*ir.Func),
	}
}

func (c *Compiler) Compile(node ast.Node) (value.Value, error) {
	fmt.Println(node)
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			_, err := c.Compile(s)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ast.ExpressionStatement:
		exp, err := c.Compile(node.Expression)
		if err != nil {
			return nil, err
		}
		return exp, nil
	case *ast.InfixExpression:
		left, err := c.Compile(node.Left)
		if err != nil {
			return nil, err
		}

		leftCon := left.(constant.Constant)

		right, err := c.Compile(node.Right)
		if err != nil {
			return nil, err
		}

		rightCon := right.(constant.Constant)

		switch node.Operator {
		case "+":
			return constant.NewAdd(leftCon, rightCon), nil
		default:
			return nil, fmt.Errorf("unknown operator %s", node.Operator)
		}
	// case *ast.PrefixExpression:
	// 	right, err := c.Compile(node.Right)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	switch node.Operator {
	// 	case "!":
	// 		c.emit(code.OpBang)
	// 	case "-":
	// 		c.emit(code.OpMinus)
	// 	default:
	// 		return fmt.Errorf("unknown operator %s", node.Operator)
	// 	}
	case *ast.IfExpression:
		condit, err := c.Compile(node.Condition)
		if err != nil {
			return nil, err
		}
		curBlock := c.CurBlock
		trBlock := c.CurBlock.Parent.NewBlock("true")
		c.CurBlock = trBlock
		var result value.Value
		result, err = c.Compile(node.Consequence)
		if err != nil {
			return nil, err
		}
		c.CurBlock = curBlock
		after := c.CurBlock.Parent.NewBlock("after")
		if node.Alternative != nil {
			faBlock := c.CurBlock.Parent.NewBlock("false")
			c.CurBlock = faBlock
			result, err = c.Compile(node.Alternative)
			if err != nil {
				return nil, err
			}
			c.CurBlock.NewCondBr(condit, trBlock, faBlock)
		} else {
			c.CurBlock.NewCondBr(condit, trBlock, after)
		}
		c.CurBlock = curBlock
		if trBlock.Term == nil {
			trBlock.NewBr(after)
			c.CurBlock = after
		}
		return result, nil
	case *ast.LetStatement:
		val, err := c.Compile(node.Value)
		if err != nil {
			return nil, err
		}
		con, ok := val.(constant.Constant)
		if !ok {
			return nil, fmt.Errorf("val not constant")
		}
		//alloca := c.CurBlock.NewAlloca(con.Type())
		//c.CurBlock.NewStore(con, alloca)
		c.Module.NewGlobalDef(node.Name.Value, con)
		return nil, nil
	// case *ast.ForLiteral:
	// 	origin := len(c.instructions)
	// 	err := c.Compile(node.Iterator)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	pos := c.emit(code.OpJumpNotTrue, 9999)

	// 	err = c.Compile(node.Body)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if c.lastInstructionIsPop() {
	// 		c.removeLastPop()
	// 	}

	// 	c.emit(code.OpJump, origin - len(c.instructions))
	// 	after := len(c.instructions)
	// 	c.changeOperand(pos, after)

	// 	if c.lastInstructionIsPop() {
	// 		c.removeLastPop()
	// 	}

	case *ast.IntegerLiteral:
		return constant.NewInt(types.I32, node.Value), nil
	case *ast.Boolean:
		return constant.NewBool(node.Value), nil
	case *ast.StringLiteral:
		return constant.NewCharArrayFromString(node.Value), nil
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			ca, err := c.Compile(s)
			if err != nil {
				return nil, err
			}
			return ca, nil
		}
	case *ast.Identifier:
		for _, s := range c.Module.Globals {
			if s.GlobalName == node.Token.Literal {
				return s, nil
			}
		}
		return nil, fmt.Errorf("could not find identifier")
	case *ast.ArrayLiteral:
		elements := make([]constant.Constant, 0)
		for _, el := range node.Elements {
			ele, err := c.Compile(el)
			if err != nil {
				return nil, err
			}
			elem := ele.(constant.Constant)
			elements = append(elements, elem)
		}
		array := constant.NewArray(elements...)
		return array, nil
	// case *ast.HashLiteral:
	// 	keys := []ast.Expression{}
	// 	for k := range node.Pairs {
	// 		keys = append(keys, k)
	// 	}
	// 	sort.Slice(keys, func(i, j int) bool {
	// 		return keys[i].String() < keys[j].String()
	// 	})

	// 	for _, k := range keys {
	// 		err := c.Compile(k)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		err = c.Compile(node.Pairs[k])
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}

	// 	c.emit(code.OpHash, len(node.Pairs)*2)
	// case *ast.IndexExpression:
	// 	exp, err := c.Compile(node.Left)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	index, err := c.Compile(node.Index)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	indexNum := index.(*constant.Int)
	// 	if exp.Type().Equal(types.NewArray(0, types.I32)) {
	// 		return c.CurBlock.NewGetElementPtr(exp, indexNum), nil
	// 	} else {
	// 		glob := exp.(*ir.Global)
	// 	}
	case *ast.FunctionLiteral:
		fun := c.Module.NewFunc(node.Name.Value, azulaTypesToLLVM[node.ReturnType.Value])
		c.Functions[node.Name.Value] = fun
		entry := fun.NewBlock("")
		c.CurBlock = entry
		ret, err := c.Compile(node.Body)
		if err != nil {
			return nil, err
		}
		c.CurBlock.NewRet(ret)
		c.CurBlock = nil
	case *ast.ReturnStatement:
		val, err := c.Compile(node.ReturnValue)
		if err != nil {
			return nil, err
		}
		return val, nil
	case *ast.CallExpression:
		f, ok := node.Function.(*ast.Identifier)
		if !ok {
			return nil, fmt.Errorf("not a function identifier")
		}
		call := c.CurBlock.NewCall(c.Functions[f.Value])
		return call, nil
	}

	return nil, nil
}

// func (c *Compiler) emit(op code.Opcode, operands ...int) int {
// 	ins := code.Make(op, operands...)
// 	pos := c.addInstruction(ins)

// 	c.setLastInstruction(op, pos)

// 	return pos
// }

// func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
// 	previous := c.scopes[c.scopeIndex].lastInstruction
// 	last := EmittedInstruction{Opcode: op, Position: pos}

// 	c.scopes[c.scopeIndex].previousInstruction = previous
// 	c.scopes[c.scopeIndex].lastInstruction = last
// }

// func (c *Compiler) addInstruction(ins []byte) int {
// 	posNewInstruction := len(c.currentInstructions())
// 	updatedInstructions := append(c.currentInstructions(), ins...)
// 	c.scopes[c.scopeIndex].instructions = updatedInstructions
// 	return posNewInstruction
// }

// func (c *Compiler) currentInstructions() code.Instructions {
// 	return c.scopes[c.scopeIndex].instructions
// }

// func (c *Compiler) Bytecode() *Bytecode {
// 	return &Bytecode{
// 		Instructions: c.currentInstructions(),
// 		Constants:    c.constants,
// 	}
// }

// func (c *Compiler) addConstant(obj object.Object) int {
// 	c.constants = append(c.constants, obj)
// 	return len(c.constants) - 1
// }

// func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
// 	if len(c.currentInstructions()) == 0 {
// 		return false
// 	}
// 	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
// }

// func (c *Compiler) removeLastPop() {
// 	last := c.scopes[c.scopeIndex].lastInstruction
// 	previous := c.scopes[c.scopeIndex].previousInstruction

// 	old := c.currentInstructions()
// 	new := old[:last.Position]

// 	c.scopes[c.scopeIndex].instructions = new
// 	c.scopes[c.scopeIndex].lastInstruction = previous
// }

// func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
// 	ins := c.currentInstructions()
// 	for i := 0; i < len(newInstruction); i++ {
// 		ins[pos+i] = newInstruction[i]
// 	}
// }

// func (c *Compiler) changeOperand(opPos int, operand int) {
// 	op := code.Opcode(c.currentInstructions()[opPos])
// 	newInstruction := code.Make(op, operand)

// 	c.replaceInstruction(opPos, newInstruction)
// }

// func (c *Compiler) enterScope() {
// 	scope := CompilationScope{
// 		instructions:        code.Instructions{},
// 		lastInstruction:     EmittedInstruction{},
// 		previousInstruction: EmittedInstruction{},
// 	}
// 	c.scopes = append(c.scopes, scope)
// 	c.scopeIndex++
// 	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
// }

// func (c *Compiler) leaveScope() code.Instructions {
// 	instructions := c.currentInstructions()

// 	c.scopes = c.scopes[:len(c.scopes)-1]
// 	c.scopeIndex--

// 	c.symbolTable = c.symbolTable.Outer

// 	return instructions
// }
