package myloxgo

import "strconv"

// Scanner は字句のスキャンを行うための構造体.java実装のloxにおけるScannerクラス.
type Scanner struct {
	source   string
	tokens   []Token
	start    int
	current  int
	line     int
	keywords map[string]TokenType
}

// NewScanner はScannerのコンストラクタ.
func NewScanner(source string) *Scanner {
	keywords := map[string]TokenType{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"for":    FOR,
		"fun":    FUN,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
	return &Scanner{
		source:   source,
		start:    0,
		current:  0,
		line:     1,
		keywords: keywords,
	}
}

// ScanTokens はスキャンのエントリーポイントとなるメソッド.
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, *NewToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '!':
		s.addToken(map[bool]TokenType{true: BANG_EQUAL, false: BANG}[s.match('=')], nil)
	case '=':
		s.addToken(map[bool]TokenType{true: EQUAL_EQUAL, false: EQUAL}[s.match('=')], nil)
	case '<':
		s.addToken(map[bool]TokenType{true: LESS_EQUAL, false: LESS}[s.match('=')], nil)
	case '>':
		s.addToken(map[bool]TokenType{true: GREATER_EQUAL, false: GREATER}[s.match('=')], nil)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			scannerError(s.line, "unexpected character")
		}
	}
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len([]rune(s.source))
}

func (s *Scanner) advance() rune {
	char := []rune(s.source)[s.current]
	s.current++
	return char
}

func (s *Scanner) addToken(typ TokenType, literal any) {
	text := string([]rune(s.source)[s.start:s.current])
	s.tokens = append(s.tokens, *NewToken(typ, text, literal, s.line))
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if []rune(s.source)[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return []rune(s.source)[s.current]
}

func (s Scanner) peekNext() rune {
	if s.current+1 >= len([]rune(s.source)) {
		return rune(0)
	}
	return []rune(s.source)[s.current+1]
}

func (s Scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s Scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		scannerError(s.line, "Unterminated string.")
		return
	}
	s.advance()

	value := string([]rune(s.source)[s.start+1 : s.current-1])
	s.addToken(STRING, value)
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	num, _ := strconv.ParseFloat(string([]rune(s.source)[s.start:s.current]), 64)
	s.addToken(NUMBER, num)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string([]rune(s.source)[s.start:s.current])
	typ, ok := s.keywords[text]
	if !ok {
		typ = IDENTIFIER
	}
	s.addToken(typ, nil)
}
