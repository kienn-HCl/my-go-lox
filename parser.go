package myloxgo

// Parser は再帰下降構文解析を行うための構造体.java実装のloxにおけるParserクラス.
type Parser struct {
	tokens  []Token
	current int
}

// NewParser はParserのコンストラクタ.
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Parse は構文解析のエントリーポイントとなるメソッド.
func (p *Parser) Parse() []Stmt {
	statements := make([]Stmt, 0, 100)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() Stmt {
	var stmt Stmt
	var ok bool
	if p.match(VAR) {
		stmt, ok = p.varDeclaration()
	} else {
		stmt, ok = p.statement()
	}
	if !ok {
		p.synchronize()
	}
	return stmt
}

func (p *Parser) statement() (Stmt, bool) {
	switch {
	case p.match(PRINT):
		return p.printStatement()
	case p.match(LEFT_BRACE):
		block, ok := p.block()
		if !ok {
			return nil, false
		}
		return NewBlock(block), true
	}

	return p.expressionStatement()
}

func (p *Parser) varDeclaration() (Stmt, bool) {
	name, ok := p.consume(IDENTIFIER, "Expect variable name.")
	if !ok {
		return nil, false
	}

	var initializer Expr
	if p.match(EQUAL) {
		initializer, ok = p.expression()
		if !ok {
			return nil, false
		}
	}

	_, ok = p.consume(SEMICOLON, "Expect ';' after expression.")
	if !ok {
		return nil, false
	}
	return NewVar(*name, initializer), true
}

func (p *Parser) printStatement() (Stmt, bool) {
	value, ok := p.expression()
	if !ok {
		return nil, false
	}
	_, ok = p.consume(SEMICOLON, "Expect ';' after expression.")
	if !ok {
		return nil, false
	}
	return NewPrint(value), true
}

func (p *Parser) expressionStatement() (Stmt, bool) {
	expr, ok := p.expression()
	if !ok {
		return nil, false
	}
	_, ok = p.consume(SEMICOLON, "Expect ';' after expression.")
	if !ok {
		return nil, false
	}
	return NewExpress(expr), true
}

func (p *Parser) block() ([]Stmt, bool) {
	statements := make([]Stmt, 0)

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	_, ok := p.consume(RIGHT_BRACE, "Expect '}' after block.")
	if !ok {
		return nil, false
	}

	return statements, true
}

func (p *Parser) expression() (Expr, bool) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, bool) {
	expr, ok := p.equality()
	if !ok {
		return nil, false
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, ok := p.assignment()
		if !ok {
			return nil, false
		}

		if v, ok := expr.(*Variable); ok {
			return NewAssign(v.Name, value), true
		}

		parserError(equals, "Invalid assignment target.")
	}

	return expr, true
}

func (p *Parser) equality() (Expr, bool) {
	expr, ok := p.comparison()
	if !ok {
		return nil, false
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := *p.previous()
		right, ok := p.comparison()
		if !ok {
			return nil, false
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, true
}

func (p *Parser) comparison() (Expr, bool) {
	expr, ok := p.term()
	if !ok {
		return nil, false
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := *p.previous()
		right, ok := p.term()
		if !ok {
			return nil, false
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, true
}

func (p *Parser) term() (Expr, bool) {
	expr, ok := p.factor()
	if !ok {
		return nil, false
	}

	for p.match(MINUS, PLUS) {
		operator := *p.previous()
		right, ok := p.factor()
		if !ok {
			return nil, false
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, true
}

func (p *Parser) factor() (Expr, bool) {
	expr, ok := p.unary()
	if !ok {
		return nil, false
	}

	for p.match(SLASH, STAR) {
		operator := *p.previous()
		right, ok := p.unary()
		if !ok {
			return nil, false
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, true
}

func (p *Parser) unary() (Expr, bool) {
	if p.match(BANG, MINUS) {
		operator := *p.previous()
		right, ok := p.unary()
		if !ok {
			return nil, false
		}
		return NewUnary(operator, right), true
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, bool) {
	switch {
	case p.match(FALSE):
		return NewLiteral(false), true
	case p.match(TRUE):
		return NewLiteral(true), true
	case p.match(NIL):
		return NewLiteral(nil), true
	case p.match(NUMBER, STRING):
		return NewLiteral(p.previous().Literal), true
	case p.match(IDENTIFIER):
		return NewVariable(*p.previous()), true
	case p.match(LEFT_PAREN):
		expr, ok := p.expression()
		if !ok {
			return nil, false
		}
		_, ok = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if !ok {
			return nil, false
		}
		return NewGrouping(expr), true
	}

	parserError(p.peek(), "Expect expression.")
	return nil, false
}

func (p *Parser) match(types ...TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Typ == typ
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Typ == EOF
}

func (p *Parser) peek() *Token {
	return &p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
}

func (p *Parser) consume(typ TokenType, message string) (*Token, bool) {
	if !p.check(typ) {
		parserError(p.peek(), message)
		return nil, false
	}

	return p.advance(), true
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Typ == SEMICOLON {
			return
		}

		switch p.peek().Typ {
		case CLASS:
			return
		case FOR:
			return
		case FUN:
			return
		case IF:
			return
		case PRINT:
			return
		case RETURN:
			return
		case VAR:
			return
		case WHILE:
			return
		}

		p.advance()
	}
}
