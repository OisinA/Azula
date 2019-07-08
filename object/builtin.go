package object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type {
	return BuiltinObj
}

func (b *Builtin) Inspect() string {
	return "builtin function"
}
