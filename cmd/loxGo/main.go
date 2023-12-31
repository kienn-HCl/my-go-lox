package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
    "my-lox-go"
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
	err = run(string(bytes))

	if err != nil {
		os.Exit(65)
	}
}

func runPrompt() {
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !reader.Scan() {
			break
		}
		err := run(reader.Text())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func run(source string) error {
    scan := myloxgo.NewScanner(source)
    for _, token := range scan.ScanTokens() {
        fmt.Println(token)
    }
	return nil
}
