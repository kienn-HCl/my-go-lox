package astPrinter

import (
	"fmt"
	"log"
	"my-go-lox"
	"strings"
)

type AstPrinter struct{}

func (a *AstPrinter) Print(expr mygolox.Expr) (str string) {
	if str, ok := expr.Accept(a).(string); ok {
		return str
	}
	log.Println("failed type assertion")
	return
}

func (a *AstPrinter) VisitBinaryExpr(binary mygolox.Binary) any {
	return a.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (a *AstPrinter) VisitGroupingExpr(grouping mygolox.Grouping) any {
	return a.parenthesize("group", grouping.Expression)
}

func (a *AstPrinter) VisitLiteralExpr(literal mygolox.Literal) any {
	if literal.Value == nil {
		return nil
	}
	return fmt.Sprint(literal.Value)
}

func (a *AstPrinter) VisitUnaryExpr(unary mygolox.Unary) any {
	return a.parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...mygolox.Expr) string {
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
