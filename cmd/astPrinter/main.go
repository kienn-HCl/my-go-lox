package main

import (
	"fmt"
	"my-go-lox"
	"my-go-lox/pkg/astPrinter"
)

func main() {
	expression := mygolox.NewBinary(
		mygolox.NewUnary(
			*mygolox.NewToken(mygolox.MINUS, "-", nil, 1),
			mygolox.NewLiteral(123)),
		*mygolox.NewToken(mygolox.STAR, "*", nil, 1),
		mygolox.NewGrouping(
			mygolox.NewLiteral(45.67)))

	printer := &astPrinter.AstPrinter{}
	fmt.Println(printer.Print(expression))
}
