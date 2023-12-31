package myloxgo

type Expr interface {
	Accept(visitor Visitor) any
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

func (e Binary) Accept(visitor Visitor) any {
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

func (e Grouping) Accept(visitor Visitor) any {
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

func (e Literal) Accept(visitor Visitor) any {
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

func (e Unary) Accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(e)
}

type Visitor interface {
	VisitBinaryExpr(binary Binary) any
	VisitGroupingExpr(grouping Grouping) any
	VisitLiteralExpr(literal Literal) any
	VisitUnaryExpr(unary Unary) any
}
