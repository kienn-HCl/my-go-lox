package myloxgo

type Expr interface {
	Accept(visitor VisitorExpr) any
}

type Assign struct {
	Name  Token
	Value Expr
}

func NewAssign(Name Token, Value Expr) *Assign {
	return &Assign{
		Name:  Name,
		Value: Value,
	}
}

func (e Assign) Accept(visitor VisitorExpr) any {
	return visitor.VisitAssignExpr(e)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func NewBinary(Left Expr, Operator Token, Right Expr) *Binary {
	return &Binary{
		Left:     Left,
		Operator: Operator,
		Right:    Right,
	}
}

func (e Binary) Accept(visitor VisitorExpr) any {
	return visitor.VisitBinaryExpr(e)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(Expression Expr) *Grouping {
	return &Grouping{
		Expression: Expression,
	}
}

func (e Grouping) Accept(visitor VisitorExpr) any {
	return visitor.VisitGroupingExpr(e)
}

type Literal struct {
	Value any
}

func NewLiteral(Value any) *Literal {
	return &Literal{
		Value: Value,
	}
}

func (e Literal) Accept(visitor VisitorExpr) any {
	return visitor.VisitLiteralExpr(e)
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func NewLogical(Left Expr, Operator Token, Right Expr) *Logical {
	return &Logical{
		Left:     Left,
		Operator: Operator,
		Right:    Right,
	}
}

func (e Logical) Accept(visitor VisitorExpr) any {
	return visitor.VisitLogicalExpr(e)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func NewUnary(Operator Token, Right Expr) *Unary {
	return &Unary{
		Operator: Operator,
		Right:    Right,
	}
}

func (e Unary) Accept(visitor VisitorExpr) any {
	return visitor.VisitUnaryExpr(e)
}

type Variable struct {
	Name Token
}

func NewVariable(Name Token) *Variable {
	return &Variable{
		Name: Name,
	}
}

func (e Variable) Accept(visitor VisitorExpr) any {
	return visitor.VisitVariableExpr(e)
}

type VisitorExpr interface {
	VisitAssignExpr(expr Assign) any
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitLogicalExpr(expr Logical) any
	VisitUnaryExpr(expr Unary) any
	VisitVariableExpr(expr Variable) any
}
