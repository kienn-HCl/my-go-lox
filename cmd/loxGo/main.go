package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"my-go-lox"
	"os"
)

var interpreter *mygolox.Interpreter = mygolox.NewInterpreter()

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

	if mygolox.HadError {
		os.Exit(65)
	}
	if mygolox.HadRuntimeError {
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
		mygolox.HadError = false
	}
}

func run(source string) {
	scan := mygolox.NewScanner(source)
	tokens := scan.ScanTokens()
	parser := mygolox.NewParser(tokens)
	statements := parser.Parse()

	if mygolox.HadError {
		return
	}

	resolver := mygolox.NewResolver(interpreter)
	resolver.ResolveStmts(statements)

	if mygolox.HadError {
		return
	}

	interpreter.Interpret(statements)
}
