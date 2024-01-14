package myloxgo

import (
	"fmt"
	"os"
)

var HadError bool = false

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
