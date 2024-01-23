package myloxgo

// Environment は環境のための構造体.java実装のloxにおけるEnvironmentクラス.
type Environment struct {
	Values map[string]any
}

// NewEnvironment はEnvironmentのコンストラクタ.
func NewEnvironment() *Environment {
	return &Environment{
		Values: map[string]any{},
	}
}

func (e *Environment) define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) get(name Token) (any, error) {
	if v, ok := e.Values[name.Lexeme]; ok {
		return v, nil
	}

	return nil, NewRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}
