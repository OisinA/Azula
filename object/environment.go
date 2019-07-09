package object

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	env := &Environment{store: s, outer: nil}
	env.Set("null", &Null{})
	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Overwrite(name string, val Object) Object {
	_, ok := e.store[name]
	if !ok && e.outer != nil {
		e.outer.Overwrite(name, val)
	} else {
		e.store[name] = val
	}
	return val
}
