package evaluator

import (
	"github.com/OisinA/Azula/ast"
	"github.com/OisinA/Azula/object"
	"github.com/OisinA/Azula/lexer"
	"github.com/OisinA/Azula/parser"
	"fmt"
	"io/ioutil"
)

var (
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL = &object.Null{}

	typeMap = map[object.ObjectType]string {
		object.INTEGER_OBJ: "int",
		object.BOOLEAN_OBJ: "bool",
		object.STRING_OBJ: "string",
		object.ARRAY_OBJ: "array",
	}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if val.Type() == object.ARRAY_OBJ {
			array := val.(*object.Array)
			if node.Name.ReturnType.Value != array.ElementType {
				return newError("trying to assign array %s to array %s: " + node.Name.Value, array.ElementType, node.Name.ReturnType.Value)
			}
			env.Set(node.Name.Value, val)
			return NULL
		}
		if val.Type() == object.CLASS_OBJ {
			class := val.(*object.Class)
			if node.Token.Literal != class.Name.String() {
				return newError("can't assign to type %s", node.Token.Literal)
			}
			env.Set(node.Name.Value, val)
			return NULL
		}
		if typeMap[val.Type()] == node.Token.Literal {
			env.Set(node.Name.Value, val)
			return NULL
		} else {
			return newError("trying to assign %s to %s: " + node.Name.Value, typeMap[val.Type()], node.Token.Literal)
		}

	case *ast.ReassignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		obj, ok := env.Get(node.Name.Value)
		if !ok {
			return newError("can't reassign value to non-existent variable '" + node.Name.Value + "'")
		}
		if typeMap[obj.Type()] != typeMap[val.Type()] {
			return newError("can't assign value of type %s to variable of type %s", typeMap[obj.Type()], typeMap[val.Type()])
		}
		env.Overwrite(node.Name.Value, val)

	case *ast.ImportStatement:
		val := Eval(node.Value, env)
		v, ok := val.(*object.String)
		if !ok {
			return newError("invalid import path")
		}
		path := v.Value
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			return newError("couldn't import file '%s'", path)
		}
		l := lexer.New(string(dat))
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			return newError("something went wrong while importing '%s'", path)
		}
		Eval(program, env)
		return NULL

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		var t string
		for _, tt := range elements {
			if t == "" {
				t = typeMap[tt.Type()]
				continue
			}
			if t != typeMap[tt.Type()] {
				return newError("trying to assign %s to array of %s", typeMap[tt.Type()], t)
			}
		}
		return &object.Array{ElementType: t, Elements: elements}

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.TypedIdentifier:
		return evalTypedIdentifier(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		function := &object.Function{Name: node.Name, Parameters: params, Env: env, Body: body, ReturnType: node.ReturnType}
		env.Set(node.Name.Token.Literal, function)
		return function

	case *ast.ClassLiteral:
		params := node.Parameters
		body := node.Body
		class := &object.Class{Name: node.Name, Parameters: params, Env: object.NewEnvironment(), Body: body}
		env.Set(node.Name.Token.Literal, class)
		return class

	case *ast.CallExpression:
		newEnv := env
		if node.Outer != nil {
			outer := node.Outer
			classEnv, ok := env.Get(outer.TokenLiteral())
			if !ok {
				return newError("couldn't find object %s", outer.TokenLiteral())
			}
			class, ok := classEnv.(*object.Class)
			if !ok {
				return newError("'%s' is not a class", node.Function.TokenLiteral())
			}
			if class.Env != nil {
				newEnv = class.Env
			} else {
				return newError("'%s' is not an object", node.Function.TokenLiteral())
			}
		}
		function := Eval(node.Function, newEnv)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, newEnv)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		switch fn := function.(type) {
		case *object.Function:
			fn = function.(*object.Function)
			result := applyFunction(function, args)
			if fn.ReturnType.Token.Literal == "void" {
				return NULL
			}
			if typeMap[result.Type()] == fn.ReturnType.Token.Literal {
				if fn.ReturnType.Token.Literal == "array" {
					array := result.(*object.Array)
					if array.ElementType != fn.ReturnType.Value {
						return newError("function %s returned array(%s), not array(%s)", fn.Name.String(), array.ElementType, fn.ReturnType.Value)
					}
				}
				return result
			} else {
				if result.Type() == object.CLASS_OBJ {
					if _, ok := env.Get(fn.ReturnType.Token.Literal); ok {
						return result
					}
				}
				return newError("function %s returned %s, not %s", fn.Name.String(), typeMap[result.Type()], fn.ReturnType.Token.Literal)
			}
		default:
			return applyFunction(function, args)
		}
	case *ast.ForLiteral:
		obj := Eval(node.Iterator, env)
		if isError(obj) {
			return obj
		}
		forLoop, ok := obj.(*object.Array)
		if !ok {
			return newError("iterator must be an array")
		}
		env1 := object.NewEnclosedEnvironment(env)
		var result object.Object
		for i := 0; i < len(forLoop.Elements); i++ {
			env1.Set(node.Parameter.String(), forLoop.Elements[i])
			result = Eval(node.Body, env1)
		}
		if result == nil {
			result = NULL
		}
		return result
	}
	return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
			return nativeBoolToBooleanObject((left.(*object.String)).Value == (right.(*object.String)).Value)
		}
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type() && operator != "+":
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case operator == "+":
		return evalStringInfixExpression(operator, &object.String{Value: string(left.Inspect())}, &object.String{Value: string(right.Inspect())})
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalTypedIdentifier(node *ast.TypedIdentifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}

	return val
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx > max {
		return newError("index out of bounds")
	}

	if idx < 0 {
		idx = int64(len(arrayObject.Elements))+idx
	}

	return arrayObject.Elements[idx]
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	case *object.Class:
		env := object.NewEnvironment()
		for paramIdx, x := range fn.Parameters {
			env.Set(x.Value, args[paramIdx])
		}
		Eval(fn.Body, env)
		return &object.Class{Name: fn.Name, Body: fn.Body, Parameters: fn.Parameters, Env: env}
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
