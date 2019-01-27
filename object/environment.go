package object

type Environment struct {
	outer   *Environment
	objects map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{objects: make(map[string]Object)}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{
		objects: make(map[string]Object),
		outer:   outer,
	}
}

func (e *Environment) Set(name string, value Object) {
	e.objects[name] = value
}
func (e *Environment) Get(name string) Object {
	value, ok := e.objects[name]
	if !ok {
		if e.outer == nil {
			return nil
		}
		return e.outer.Get(name)
	}
	return value
}
