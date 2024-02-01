package mygolox

import "fmt"

// Token は文字をスキャンして得た字句の情報を保存する構造体.
type Token struct {
	Typ     TokenType
	Lexeme  string
	Literal any
	Line    int
}

// NewToken はTokenのコンストラクタ.
func NewToken(typ TokenType, lexeme string, literal any, line int) *Token {
	return &Token{
		Typ:     typ,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %v", t.Typ, t.Lexeme, t.Literal)
}
