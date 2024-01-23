package myloxgo

import (
	"fmt"
	"os"
)

// HadError は字句解析または構文解析の処理でエラーがあったことを伝えるフラグ.
var HadError = false

// HadRuntimeError はコードを実行する際の処理でエラーがあったことを伝えるフラグ.
var HadRuntimeError = false

func scannerError(line int, message string) {
	report(line, "", message)
}

func parserError(token *Token, message string) {
	if token.Typ == EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func report(line int, where, message string) {
	fmt.Fprintln(os.Stderr, "[line", line, "] Error", where, ":", message)
	HadError = true
}

// RuntimeError はランタイムエラーを報告するための構造体.errorインターフェイスを満たす.
type RuntimeError struct {
	Token   Token
	Message string
}

// NewRuntimeError はRuntimeErrorのコンストラクタ.
func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{
		Token:   token,
		Message: message,
	}
}

func (r *RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %d]", r.Message, r.Token.Line)
}
