package main

import (
	"fmt"
	"my-lox-go"
	"my-lox-go/pkg/astPrinter"
)

func main() {
	expression := myloxgo.NewBinary(
		myloxgo.NewUnary(
			*myloxgo.NewToken(myloxgo.MINUS, "-", nil, 1),
			myloxgo.NewLiteral(123)),
		*myloxgo.NewToken(myloxgo.STAR, "*", nil, 1),
		myloxgo.NewGrouping(
			myloxgo.NewLiteral(45.67)))

	printer := &astPrinter.AstPrinter{}
	fmt.Println(printer.Print(expression))
}
