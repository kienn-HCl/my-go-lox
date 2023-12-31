package main

import (
	"fmt"
	"log"
	"my-lox-go"
	"strings"
)

type AstPrinter struct{}

func (a *AstPrinter) print(expr myloxgo.Expr) string {
	if str, ok := expr.Accept(a).(string); ok {
		return str
	}
	log.Println("failed type assertion")
	return ""
}

func (a *AstPrinter) VisitBinaryExpr(binary myloxgo.Binary) any {
	return a.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (a *AstPrinter) VisitGroupingExpr(grouping myloxgo.Grouping) any {
	return a.parenthesize("group", grouping.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(literal myloxgo.Literal) any {
	if literal.Value == nil {
		return nil
	}
	return fmt.Sprint(literal.Value)
}

func (a *AstPrinter) VisitUnaryExpr(unary myloxgo.Unary) any {
	return a.parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...myloxgo.Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		if str, ok := expr.Accept(a).(string); ok {
			builder.WriteString(str)
		}
	}
	builder.WriteString(")")

	return builder.String()
}

func main() {
	expression := myloxgo.NewBinary(
		myloxgo.NewUnary(
			*myloxgo.NewToken(myloxgo.MINUS, "-", nil, 1),
			myloxgo.NewLiteral(123)),
		*myloxgo.NewToken(myloxgo.STAR, "*", nil, 1),
		myloxgo.NewGrouping(
			myloxgo.NewLiteral(45.67)))

    printer := &AstPrinter{}
    fmt.Println(printer.print(expression))
}
