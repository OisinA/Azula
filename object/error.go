package object

type Error struct {
	Message string
}

func (e *Error) Type() Type {
	return ErrorObj
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
