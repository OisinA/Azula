package compiler

import (
	"azula/ast"

	"fmt"
)

type Type string

const (
	Int     Type = "int"
	String  Type = "string"
	Boolean Type = "bool"
	Array   Type = "array"
	Map     Type = "map"
	Void    Type = "void"
)

type Typechecker struct {
	variableTypes map[string]Type
}

func NewTypechecker() *Typechecker {
	return &Typechecker{make(map[string]Type)}
}

func NewTypecheckerFromVars(vars map[string]Type) *Typechecker {
	return &Typechecker{vars}
}

func (t *Typechecker) Typecheck(node ast.Node) (Type, error) {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			ret, err := t.Typecheck(s)
			if err != nil {
				return "", err
			}
			if ret == Void {
				continue
			}
			return ret, err
		}
		return Void, nil
	case *ast.ExpressionStatement:
		return t.Typecheck(node.Expression)
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			ret, err := t.Typecheck(s)
			if err != nil {
				return "", err
			}
			if ret == Void {
				continue
			}
			return ret, nil
		}
		return Void, nil
	case *ast.ReturnStatement:
		return t.Typecheck(node.ReturnValue)
	case *ast.LetStatement:
		retType, err := t.Typecheck(node.Value)
		if err != nil {
			return "", err
		}
		if string(retType) != node.Token.Literal {
			return "", fmt.Errorf("could not assign %s to variable of type %s", retType, node.Token.Literal)
		}

		if retType == Array {
			rType, err := t.Typecheck(node.Value.(*ast.ArrayLiteral).Elements[0])
			if err != nil {
				return "", fmt.Errorf("can't read type of array")
			}
			if node.Name.ReturnType.Value != string(rType) {
				return "", fmt.Errorf("could not assign array of type %s to array of type %s", rType, node.Name.ReturnType.Value)
			}
		} else if retType == Map {
			var keyType Type
			var valType Type
			for s := range node.Value.(*ast.HashLiteral).Pairs {
				keyType, err = t.Typecheck(s)
				if err != nil {
					return "", err
				}
				valType, err = t.Typecheck(node.Value.(*ast.HashLiteral).Pairs[s])
				if err != nil {
					return "", err
				}
				break
			}
			if node.Name.ReturnType.Value != string(keyType) {
				return "", fmt.Errorf("could not assign map with key %s to map (%s=>%s)", node.Name.ReturnType.Value, keyType, valType)
			}
		}
		t.variableTypes[node.Name.Value] = retType
		return Void, nil
	case *ast.InfixExpression:
		left, err := t.Typecheck(node.Left)
		right, err := t.Typecheck(node.Right)
		if err != nil {
			return Void, err
		}
		if left != right {
			return "", fmt.Errorf("could not perform operation between %s and %s", left, right)
		}
		return left, nil
	case *ast.PrefixExpression:
		right, err := t.Typecheck(node.Right)
		if err != nil {
			return Void, err
		}
		switch node.Operator {
		case "-":
			return Int, nil
		case "!":
			return Boolean, nil
		}
		return right, nil
	case *ast.IfExpression:
		conType, err := t.Typecheck(node.Consequence)
		if err != nil {
			return "", err
		}
		alternateType, err := t.Typecheck(node.Alternative)
		if err != nil {
			return "", err
		}
		if conType != alternateType {
			return "", fmt.Errorf("conditional statement can't return different types")
		}
		return conType, nil

	case *ast.IntegerLiteral:
		return Int, nil
	case *ast.StringLiteral:
		return String, nil
	case *ast.Boolean:
		return Boolean, nil
	case *ast.Identifier:
		retType, ok := t.variableTypes[node.Value]
		if !ok {
			return "", fmt.Errorf("unknown variable %s, can't read type", node.Value)
		}
		return retType, nil
	case *ast.ArrayLiteral:
		var arrayType Type
		for _, s := range node.Elements {
			ty, err := t.Typecheck(s)
			if err != nil {
				return Void, err
			}
			if arrayType == "" {
				arrayType = ty
				continue
			}
			if ty != arrayType {
				return "", fmt.Errorf("could not assign %s to array of type %s", ty, arrayType)
			}
		}
		return Array, nil
	case *ast.HashLiteral:
		var keyType Type
		var valType Type
		for s := range node.Pairs {
			key, err := t.Typecheck(s)
			val, err := t.Typecheck(node.Pairs[s])
			if err != nil {
				return "", err
			}
			if keyType == "" && valType == "" {
				keyType = key
				valType = val
				continue
			}
			if key != keyType {
				return "", fmt.Errorf("could not assign key of type %s to map of type (%s=>%s)", key, keyType, valType)
			}
			if val != valType {
				return "", fmt.Errorf("could not assign value of type %s to map of type (%s=>%s)", val, keyType, valType)
			}
		}
		return Map, nil
	}
	return Void, nil
}
