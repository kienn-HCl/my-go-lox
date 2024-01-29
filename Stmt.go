package myloxgo

type Stmt interface {
	Accept(visitor VisitorStmt) any
}

type Block struct {
	Statements []Stmt
}

func NewBlock(Statements []Stmt) *Block {
	return &Block{
		Statements: Statements,
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

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(Condition Expr, ThenBranch Stmt, ElseBranch Stmt) *If {
	return &If{
		Condition:  Condition,
		ThenBranch: ThenBranch,
		ElseBranch: ElseBranch,
	}
}

func (e If) Accept(visitor VisitorStmt) any {
	return visitor.VisitIfStmt(e)
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

type While struct {
	condition Expr
	body      Stmt
}

func NewWhile(condition Expr, body Stmt) *While {
	return &While{
		condition: condition,
		body:      body,
	}
}

func (e While) Accept(visitor VisitorStmt) any {
	return visitor.VisitWhileStmt(e)
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
	VisitIfStmt(stmt If) any
	VisitPrintStmt(stmt Print) any
	VisitWhileStmt(stmt While) any
	VisitVarStmt(stmt Var) any
}
