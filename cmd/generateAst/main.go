package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Usage: generateAst <output directory>")
		os.Exit(64)
	}
	outputDir := flag.Arg(0)
	err := defineAst(outputDir, "Expr", []string{
		"Binary     : Left Expr, Operator Token, Right Expr",
		"Grouping   : Expression Expr",
		"Literal    : Value any",
		"Unary      : Operator Token, Right Expr",
	})
	if err != nil {
		log.Fatalln(err)
	}
	err = defineAst(outputDir, "Stmt", []string{
		"Express    : Expression Expr",
		"Print      : Expression Expr",
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func defineAst(outputDir, baseName string, types []string) (err error) {
	path := filepath.Join(outputDir, baseName+".go")
	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err = writer.Close()
	}()

	fmt.Fprintln(writer, "package myloxgo")
	fmt.Fprintln(writer)
	// fmt.Fprintln(writer, "import")
	// fmt.Fprintln(writer)

	// define baseName interface
	fmt.Fprintln(writer, "type", baseName, "interface {")
	fmt.Fprintln(writer, "	Accept(visitor Visitor"+baseName+") any")
	fmt.Fprintln(writer, "}")
	fmt.Fprintln(writer)

	// AST structs
	for _, typ := range types {
		typSplit := strings.Split(typ, ":")
		structName := strings.TrimSpace(typSplit[0])
		fields := strings.TrimSpace(typSplit[1])
		defineType(writer, baseName, structName, fields)
	}

	defineVisitor(writer, baseName, types)

	return
}

func defineType(writer io.Writer, baseName, structName, fieldList string) {
	// define struct
	fmt.Fprintln(writer, "type", structName, "struct {")
	fields := strings.Split(fieldList, ", ")
	for _, field := range fields {
		fmt.Fprintln(writer, "	"+field)
	}
	fmt.Fprintln(writer, "}")
	fmt.Fprintln(writer)

	// constructer
	fmt.Fprintln(writer, "func New"+structName+"("+fieldList+")", "*"+structName, "{")
	fmt.Fprintln(writer, "	return &"+structName+"{")
	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		fmt.Fprintln(writer, "		"+name+":", name+",")
	}
	fmt.Fprintln(writer, "	}")
	fmt.Fprintln(writer, "}")
	fmt.Fprintln(writer)

	// define accept
	fmt.Fprintln(writer, "func (e", structName+")", "Accept(visitor Visitor"+baseName+")", "any", "{")
	fmt.Fprintln(writer, "	return visitor.Visit"+structName+baseName+"(e)")
	fmt.Fprintln(writer, "}")

	fmt.Fprintln(writer)
}

func defineVisitor(writer io.Writer, baseName string, types []string) {
	fmt.Fprintln(writer, "type", "Visitor"+baseName, "interface {")
	for _, typ := range types {
		splitedTyp := strings.Split(typ, ":")[0]
		typeName := strings.TrimSpace(splitedTyp)
		fmt.Fprintln(writer, "	Visit"+typeName+baseName+"("+strings.ToLower(baseName), typeName+")", "any")
	}
	fmt.Fprintln(writer, "}")
	fmt.Fprintln(writer)
}
