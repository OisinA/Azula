package object

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	STRING_OBJ = "STRING"
	ARRAY_OBJ = "ARRAY"
	BUILTIN_OBJ = "BUILTIN"
	FOR_OBJ = "FOR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type BuiltinFunction func(args ...Object) Object
