package glox

import "fmt"

type Enviorment struct {
	enclosing *Enviorment
	values    map[string]any
}

func (e *Enviorment) Get(name Token) (any, error) {
	v, ok := e.values[name.Lexme]
	if ok {
		return v, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, RunTimeError{t: name, m: fmt.Sprintf("Undefined Variable '%v'.", name.Lexme)}
}

func (e *Enviorment) Put(name string, value any) {
	e.values[name] = value
}

func (e *Enviorment) Assign(name Token, value any) error {
	_, exists := e.values[name.Lexme]
	if exists {
		e.Put(name.Lexme, value)
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return fmt.Errorf("undefined variable '%v'", name.Lexme)
}
