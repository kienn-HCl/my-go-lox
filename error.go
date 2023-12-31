package myloxgo

import (
	"fmt"
	"os"
)

var HadError bool = false

func error(line int, message string) {
    report(line, "", message)
}

func report(line int, where ,message string) {
    fmt.Fprintln(os.Stderr, "[line", line, "] Error", where, ":", message)
    HadError = true
}
