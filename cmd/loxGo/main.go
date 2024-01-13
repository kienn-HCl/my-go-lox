package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"my-lox-go"
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
	}
}

func run(source string) {
	scan := myloxgo.NewScanner(source)
	for _, token := range scan.ScanTokens() {
		fmt.Println(token)
	}
}
