package myloxgo

import (
	"fmt"
	"os"
)

// Interpreter は構文木を解釈するための構造体.java実装のloxにおけるInterpreterクラス.
type Interpreter struct {
	Globals     *Environment
	Environment *Environment
}

// NewInterpreter はInterpreterのコンストラクタ.
func NewInterpreter() *Interpreter {
	global := NewEnvironment()
	global.define("clock", NewClock())
	return &Interpreter{
		Globals:     global,
		Environment: global,
	}
}

// Interpret は構文木を実行するためのエントリーポイントとなるメソッド.
func (i *Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			HadRuntimeError = true
			return
		}
	}
}

func (i *Interpreter) VisitBinaryExpr(expr Binary) any {
	left := i.evaluate(expr.Left)
	if v, ok := left.(error); ok {
		return v
	}
	right := i.evaluate(expr.Right)
	if v, ok := right.(error); ok {
		return v
	}

	switch expr.Operator.Typ {
	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)
	case GREATER:
		value := i.checkNumberOperands(expr.Operator, left, right, func(a1, a2 float64) any { return a1 > a2 })
		return value
	case GREATER_EQUAL:
		value := i.checkNumberOperands(expr.Operator, left, right, func(a1, a2 float64) any { return a1 >= a2 })
		return value
	case LESS:
		value := i.checkNumberOperands(expr.Operator, left, right, func(a1, a2 float64) any { return a1 < a2 })
		return value
	case LESS_EQUAL:
		value := i.checkNumberOperands(expr.Operator, left, right, func(a1, a2 float64) any { return a1 <= a2 })
		return value
	case MINUS:
		value := i.checkNumberOperands(expr.Operator, left, right, func(a1, a2 float64) any { return a1 - a2 })
		return value
	case PLUS:
		value := i.checkNumberOperands(expr.Operator, left, right, func(a1, a2 float64) any { return a1 + a2 })
		if v, ok := value.(float64); ok {
			return v
		}
		value = i.checkStringOperands(expr.Operator, left, right, func(a1, a2 string) any { return a1 + a2 })
		if v, ok := value.(string); ok {
			return v
		}

		return NewRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
	case SLASH:
		value := i.checkNumberOperands(expr.Operator, left, right, func(a1, a2 float64) any { return a1 / a2 })
		return value
	case STAR:
		value := i.checkNumberOperands(expr.Operator, left, right, func(a1, a2 float64) any { return a1 * a2 })
		return value
	}

	return nil
}

func (i *Interpreter) VisitCallExpr(expr Call) any {
	callee := i.evaluate(expr.Callee)
	if v, ok := callee.(error); ok {
		return v
	}

	arguments := make([]any, 0)
	for _, argument := range expr.Arguments {
		arg := i.evaluate(argument)
		if v, ok := arg.(error); ok {
			return v
		}
		arguments = append(arguments, arg)
	}

	if function, ok := callee.(LoxCallable); ok {
		if len(arguments) != function.Arity() {
			return NewRuntimeError(expr.Paren, "Expected "+fmt.Sprint(function.Arity())+" arguments but got "+fmt.Sprint(len(arguments))+".")
		}
		return function.Call(*i, arguments)
	}
	return NewRuntimeError(expr.Paren, "Can only call functions and classes.")
}

func (i *Interpreter) VisitUnaryExpr(expr Unary) any {
	right := i.evaluate(expr.Right)
	if v, ok := right.(error); ok {
		return v
	}

	switch expr.Operator.Typ {
	case BANG:
		return !i.isTruthy(right)
	case MINUS:
		value := i.checkNumberOperand(expr.Operator, right, func(f float64) any { return -f })
		return value
	}

	return nil
}

func (i *Interpreter) VisitGroupingExpr(expr Grouping) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr Literal) any {
	return expr.Value
}

func (i *Interpreter) VisitLogicalExpr(expr Logical) any {
	left := i.evaluate(expr.Left)
	if err, ok := left.(error); ok {
		return err
	}

	if expr.Operator.Typ == OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}
	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitVariableExpr(expr Variable) any {
	value, err := i.Environment.get(expr.Name)
	if err != nil {
		return err
	}
	return value
}

func (i *Interpreter) VisitAssignExpr(expr Assign) any {
	value := i.evaluate(expr.Value)
	if err, ok := value.(error); ok {
		return err
	}

	err := i.Environment.assign(expr.Name, value)
	if err != nil {
		return err
	}
	return value
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) any {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) any {
	previous := i.Environment
	defer func() {
		i.Environment = previous
	}()

	i.Environment = environment

	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt Block) any {
	err := i.executeBlock(stmt.Statements, NewEnvironment().ChangeEnclosing(i.Environment))
	if err != nil {
		return err
	}
	return nil
}

func (i *Interpreter) VisitExpressStmt(stmt Express) any {
	value := i.evaluate(stmt.Expression)
	if err, ok := value.(error); ok {
		return err
	}
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt Function) any {
	function := NewLoxFunction(&stmt, i.Environment)
	i.Environment.define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt If) any {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		err := i.execute(stmt.ThenBranch)
		if err != nil {
			return err
		}
	} else if stmt.ElseBranch != nil {
		err := i.execute(stmt.ElseBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt Print) any {
	value := i.evaluate(stmt.Expression)
	if err, ok := value.(error); ok {
		return err
	}
	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt Return) any {
	var value any = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
		if err, ok := value.(error); ok {
			return err
		}
	}

	return NewReturnValue(value)
}

func (i *Interpreter) VisitVarStmt(stmt Var) any {
	var value any
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
		if err, ok := value.(error); ok {
			return err
		}
	}

	i.Environment.define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt While) any {
	for i.isTruthy(i.evaluate(stmt.condition)) {
		i.execute(stmt.body)
	}

	return nil
}

func (i *Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}
	if v, ok := object.(bool); ok {
		return v
	}
	return true
}
func (i *Interpreter) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) stringify(object any) any {
	if object == nil {
		return "nil"
	}
	return object
}

func (i *Interpreter) checkNumberOperand(operator Token, operand any, calc func(float64) any) any {
	if v, ok := operand.(float64); ok {
		return calc(v)
	}

	return NewRuntimeError(operator, "Operand must be a number.")
}

func (i *Interpreter) checkNumberOperands(operator Token, left, right any, calc func(float64, float64) any) any {
	vl, okl := left.(float64)
	vr, okr := right.(float64)
	if okl && okr {
		return calc(vl, vr)
	}

	return NewRuntimeError(operator, "Operand must be a numbers.")
}

func (i *Interpreter) checkStringOperands(operator Token, left, right any, calc func(string, string) any) any {
	vl, okl := left.(string)
	vr, okr := right.(string)
	if okl && okr {
		return calc(vl, vr)
	}

	return NewRuntimeError(operator, "Operand must be a string.")
}
