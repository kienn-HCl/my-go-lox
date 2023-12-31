//go:generate stringer -type=TokenType
package myloxgo

type TokenType int

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
