package object

type Environment struct {
	objects map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{objects: make(map[string]Object)}
}

func (e *Environment) Set(name string, value Object) {
	e.objects[name] = value
}
func (e *Environment) Get(name string) Object {
	value, ok := e.objects[name]
	if !ok {
		return nil
	}
	return value
}
