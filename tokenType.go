// Package myloxgo はloxインタープリタのgo実装.
//go:generate stringer -type=TokenType
package myloxgo

// TokenType は字句の種類.
type TokenType int

// TokenType はenumとしてconst+iotaで実現.
const (
    // 記号1個のトークン
    LEFT_PAREN TokenType = iota + 1
    RIGHT_PAREN
    LEFT_BRACE
    RIGHT_BRACE
    COMMA
    DOT
    MINUS
    PLUS
    SEMICOLON
    SLASH
    STAR

    // 記号1個または2個によるトークン
    BANG
    BANG_EQUAL
    EQUAL
    EQUAL_EQUAL
    GREATER
    GREATER_EQUAL
    LESS
    LESS_EQUAL

    // リテラル
    IDENTIFIER
    STRING
    NUMBER

    // キーワード
    AND
    CLASS
    ELSE
    FALSE
    FUN
    FOR
    IF
    NIL
    OR
    PRINT
    RETURN
    SUPER
    THIS
    TRUE
    VAR
    WHILE

    EOF
)
