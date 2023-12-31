package myloxgo

type Expr interface {
	accept(visitor Visitor) any
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

func (e Binary) accept(visitor Visitor) any {
	return visitor.visitBinaryExpr(e)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(Expression Expr) *Grouping {
	return &Grouping{
		Expression: Expression,
	}
}

func (e Grouping) accept(visitor Visitor) any {
	return visitor.visitGroupingExpr(e)
}

type Literal struct {
	Value any
}

func NewLiteral(Value any) *Literal {
	return &Literal{
		Value: Value,
	}
}

func (e Literal) accept(visitor Visitor) any {
	return visitor.visitLiteralExpr(e)
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

func (e Unary) accept(visitor Visitor) any {
	return visitor.visitUnaryExpr(e)
}

type Visitor interface {
	visitBinaryExpr(binary Binary) any
	visitGroupingExpr(grouping Grouping) any
	visitLiteralExpr(literal Literal) any
	visitUnaryExpr(unary Unary) any
}
