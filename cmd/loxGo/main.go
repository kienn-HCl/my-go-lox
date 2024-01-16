package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"my-lox-go"
	"os"
)

var interpreter *myloxgo.Interpreter = &myloxgo.Interpreter{}

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

	if myloxgo.HadError {
		os.Exit(65)
	}
	if myloxgo.HadRuntimeError {
		os.Exit(70)
	}
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
	statements := parser.Parse()

	if myloxgo.HadError {
		return
	}

	interpreter.Interpret(statements)
}
