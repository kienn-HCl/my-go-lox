package myloxgo

type Stmt interface {
	Accept(visitor VisitorStmt) any
}

type Block struct {
	statements []Stmt
}

func NewBlock(statements []Stmt) *Block {
	return &Block{
		statements: statements,
	}
}

func (e Block) Accept(visitor VisitorStmt) any {
	return visitor.VisitBlockStmt(e)
}

type Express struct {
	Expression Expr
}

func NewExpress(Expression Expr) *Express {
	return &Express{
		Expression: Expression,
	}
}

func (e Express) Accept(visitor VisitorStmt) any {
	return visitor.VisitExpressStmt(e)
}

type Print struct {
	Expression Expr
}

func NewPrint(Expression Expr) *Print {
	return &Print{
		Expression: Expression,
	}
}

func (e Print) Accept(visitor VisitorStmt) any {
	return visitor.VisitPrintStmt(e)
}

type Var struct {
	Name        Token
	Initializer Expr
}

func NewVar(Name Token, Initializer Expr) *Var {
	return &Var{
		Name:        Name,
		Initializer: Initializer,
	}
}

func (e Var) Accept(visitor VisitorStmt) any {
	return visitor.VisitVarStmt(e)
}

type VisitorStmt interface {
	VisitBlockStmt(stmt Block) any
	VisitExpressStmt(stmt Express) any
	VisitPrintStmt(stmt Print) any
	VisitVarStmt(stmt Var) any
}
