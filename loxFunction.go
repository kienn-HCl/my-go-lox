package mygolox

type LoxFunction struct {
	declaration *Function
	closure     *Environment
}

func NewLoxFunction(declaration *Function, closure *Environment) *LoxFunction {
	return &LoxFunction{
		declaration: declaration,
		closure:     closure,
	}
}

func (l *LoxFunction) Call(interpreter Interpreter, arguments []any) any {
	environment := NewEnvironment().ChangeEnclosing(l.closure)
	for i, param := range l.declaration.Params {
		environment.define(param.Lexeme, arguments[i])
	}

	ret := interpreter.executeBlock(l.declaration.Body, environment)
	if err, ok := ret.(error); ok {
		return err
	}
	if r, ok := ret.(*ReturnValue); ok {
		return r.value
	}
	return nil
}

func (l *LoxFunction) Arity() int {
	return len(l.declaration.Params)
}

func (l *LoxFunction) String() string {
	return "<fn " + l.declaration.Name.Lexeme + ">"
}
