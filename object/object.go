package object

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	ARRAY_OBJ        = "ARRAY"
	BUILTIN_OBJ      = "BUILTIN"
	FOR_OBJ          = "FOR"
	CLASS_OBJ        = "CLASS"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

func Equality(obj1 *Object, obj2 *Object) bool {
	if (*obj1).Type() != (*obj2).Type() {
		return false
	}
	switch (*obj1).Type() {
	case INTEGER_OBJ:
		int1 := ((*obj1).(*Integer))
		int2 := ((*obj2).(*Integer))
		return int1.Value == int2.Value
	case STRING_OBJ:
		str1 := ((*obj1).(*String))
		str2 := ((*obj2).(*String))
		return str1.Value == str2.Value
	default:
		return obj1 == obj2
	}
}

type BuiltinFunction func(args ...Object) Object
