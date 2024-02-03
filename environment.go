package mygolox

// Environment は環境のための構造体.java実装のloxにおけるEnvironmentクラス.
type Environment struct {
	Enclosing *Environment
	Values    map[string]any
}

// NewEnvironment はEnvironmentのコンストラクタ.
func NewEnvironment() *Environment {
	return &Environment{
		Enclosing: nil,
		Values:    map[string]any{},
	}
}

// ChangeEnclosing はEnvironment構造体のEnclosingフィールドを変更する.
// コンストラクタとともに使われるのを想定している.
// ex) NewEnvironment().ChangeEnclosing(enclosing)
func (e *Environment) ChangeEnclosing(enclosing *Environment) *Environment {
	e.Enclosing = enclosing
	return e
}

func (e *Environment) define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.Enclosing
	}

	return environment
}

func (e *Environment) getAt(distance int, name string) any {
	return e.ancestor(distance).Values[name]
}

func (e *Environment) assignAt(distance int, name Token, value any) {
	e.ancestor(distance).Values[name.Lexeme] = value
}

func (e *Environment) get(name Token) (any, error) {
	if v, ok := e.Values[name.Lexeme]; ok {
		return v, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.get(name)
	}

	return nil, NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}

func (e *Environment) assign(name Token, value any) error {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.assign(name, value)
	}

	return NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}
