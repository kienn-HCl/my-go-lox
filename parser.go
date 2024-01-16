package myloxgo

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() []Stmt {
	statements := make([]Stmt, 0, 100)
	for !p.isAtEnd() {
		stmt := p.statement()
		if stmt == nil {
			p.synchronize()
		}
		statements = append(statements, stmt)
	}

	return statements
}

func (p *Parser) expression() (Expr, bool) {
	return p.equality()
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value, ok := p.expression()
	if !ok {
		return nil
	}
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return NewPrint(value)
}

func (p *Parser) expressionStatement() Stmt {
	expr, ok := p.expression()
	if !ok {
		return nil
	}
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return NewExpress(expr)
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
