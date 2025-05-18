package glox

import "fmt"

type Enviorment struct {
	values map[string]any
}

func (e *Enviorment) Get(name Token) (any, error) {
	v, ok := e.values[name.Lexme]
	if !ok {
		return nil, RunTimeError{t: name, m: fmt.Sprintf("Undefined Variable '%v'.", name.Lexme)}
	}
	return v, nil
}

func (e *Enviorment) Put(name string, value any) {
	e.values[name] = value
}

func (e *Enviorment) Assign(name Token, value any) error {
	_, exists := e.values[name.Lexme]
	if !exists {
		return fmt.Errorf("Undefined variable '%v'.", name.Lexme)
	}
	e.Put(name.Lexme, value)
	return nil
}
