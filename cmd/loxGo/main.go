package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"my-lox-go"
	"my-lox-go/pkg/astPrinter"
	"os"
)

func main() {
	flag.Parse()
	if flag.NArg() > 1 {
		fmt.Println("Usage: lox-go [script]")
		os.Exit(64)
	} else if flag.NArg() == 1 {
		runFile(flag.Arg(0))
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	run(string(bytes))
}

func runPrompt() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !reader.Scan() {
			break
		}
		run(reader.Text())
		myloxgo.HadError = false
	}
}

func run(source string) {
	scan := myloxgo.NewScanner(source)
	tokens := scan.ScanTokens()
	parser := myloxgo.NewParser(tokens)
	expression := parser.Parse()

	if myloxgo.HadError {
		return
	}

	printer := &astPrinter.AstPrinter{}
	fmt.Println(printer.Print(expression))
}
