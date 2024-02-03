package mygolox

// Resolver は変数解決のための構造体.java実装のloxにおけるResolverクラス.
type Resolver struct {
	Interpreter     *Interpreter
	Scopes          *stack
	currentFunction FunctionType
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		Interpreter:     interpreter,
		Scopes:          newStack(),
		currentFunction: NONE,
	}
}

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
)

type stack struct {
	dataGroup []map[string]bool
}

func newStack() *stack {
	return &stack{
		dataGroup: []map[string]bool{},
	}
}

func (s *stack) size() int {
	return len(s.dataGroup)
}

func (s *stack) push(data map[string]bool) {
	s.dataGroup = append(s.dataGroup, data)
}

func (s *stack) pop() map[string]bool {
	value := s.dataGroup[s.size()-1]
	s.dataGroup = s.dataGroup[:s.size()-1]
	return value
}

func (s *stack) isEmpty() bool {
	return s.size() == 0
}

func (s *stack) peek() *map[string]bool {
	if s.isEmpty() {
		return nil
	}
	return &s.dataGroup[s.size()-1]
}

func (s *stack) get(i int) *map[string]bool {
	if i < 0 || s.size() <= i {
		return nil
	}

	return &s.dataGroup[i]
}

func (r *Resolver) ResolveStmts(statements []Stmt) {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
}

func (r *Resolver) VisitBlockStmt(stmt Block) any {
	r.beginScope()
	r.ResolveStmts(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitExpressStmt(stmt Express) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt If) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt Print) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt Return) any {
	if r.currentFunction == NONE {
		parserResolverError(&stmt.Keyword, "Can't return from top-level code.")
	}
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}

	return nil
}

func (r *Resolver) VisitWhileStmt(stmt While) any {
	r.resolveExpr(stmt.condition)
	r.resolveStmt(stmt.body)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt Function) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) VisitVarStmt(stmt Var) any {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr Assign) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr Binary) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr Call) any {
	r.resolveExpr(expr.Callee)

	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}

	return nil
}

func (r *Resolver) VisitGroupingExpr(expr Grouping) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr Literal) any {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr Logical) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr Variable) any {
	// 変数がそれ自身の初期化子の中でアクセスされてるかのチェック(ex: var a = a;).
	if !r.Scopes.isEmpty() {
		if v, ok := (*r.Scopes.peek())[expr.Name.Lexeme]; ok && !v {
			parserResolverError(&expr.Name, "Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveFunction(function Function, typ FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = typ
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStmts(function.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) beginScope() {
	r.Scopes.push(map[string]bool{})
}

func (r *Resolver) endScope() {
	r.Scopes.pop()
}

func (r *Resolver) declare(name Token) {
	if r.Scopes.isEmpty() {
		return
	}
	scope := *r.Scopes.peek()
	if _, ok := scope[name.Lexeme]; ok {
		parserResolverError(&name, "Already a variable with this name in this scope.")
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name Token) {
	if r.Scopes.isEmpty() {
		return
	}
	(*r.Scopes.peek())[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := r.Scopes.size() - 1; i >= 0; i-- {
		if _, ok := (*r.Scopes.get(i))[name.Lexeme]; ok {
			r.Interpreter.resolve(expr, r.Scopes.size()-1-i)
			return
		}
	}
}
