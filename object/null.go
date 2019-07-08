package object

type Null struct {

}

func (n *Null) Type() Type {
	return NullObj
}

func (n *Null) Inspect() string {
	return "null"
}
