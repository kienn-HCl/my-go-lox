package myloxgo

type Expr interface {
	Accept(visitor VisitorExpr) any
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

type VisitorExpr interface {
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitUnaryExpr(expr Unary) any
}
