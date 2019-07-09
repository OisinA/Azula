package object

// Type is type of Object as String
type Type string

const (
	// IntegerObj - int
	IntegerObj = "INTEGER"
	// BooleanObj - bool
	BooleanObj = "BOOLEAN"
	// NullObj - null
	NullObj = "NULL"
	// ReturnValueObj - return value
	ReturnValueObj = "RETURN_VALUE"
	// ErrorObj - error
	ErrorObj = "ERROR"
	// FunctionObj - function
	FunctionObj = "FUNCTION"
	// StringObj - string
	StringObj = "STRING"
	// ArrayObj - array
	ArrayObj = "ARRAY"
	// BuiltinObj - builtin function
	BuiltinObj = "BUILTIN"
	// ForObj - for
	ForObj = "FOR"
	// ClassObj - class
	ClassObj = "CLASS"
	// HashObj - maps
	HashObj = "HASH"
)

// Object is azula object
type Object interface {
	Type() Type
	Inspect() string
}

// Equality returns if two objects are equal
func Equality(obj1 *Object, obj2 *Object) bool {
	if (*obj1).Type() != (*obj2).Type() {
		return false
	}
	switch (*obj1).Type() {
	case IntegerObj:
		int1 := ((*obj1).(*Integer))
		int2 := ((*obj2).(*Integer))
		return int1.Value == int2.Value
	case StringObj:
		str1 := ((*obj1).(*String))
		str2 := ((*obj2).(*String))
		return str1.Value == str2.Value
	case NullObj:
		_, ok := (*obj2).(*Null)
		if ok {
			return true
		}
		return false
	default:
		return obj1 == obj2
	}
}

type BuiltinFunction func(args ...Object) Object

type HashKey struct {
	Type  Type
	Value uint64
}
