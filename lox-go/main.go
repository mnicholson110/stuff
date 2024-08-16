package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mnicholson110/lox-go/lox"
)

var hadError bool = false

func main() {
	//	expression := lox.NewBinary(
	//		lox.NewUnary(
	//			lox.Token{lox.MINUS, "-", nil, 1},
	//			lox.NewLiteral(123)),
	//		lox.Token{lox.STAR, "*", nil, 1},
	//		lox.NewGrouping(lox.NewLiteral(45.67)))
	//
	//	fmt.Println(expression.AstPrint())

	switch len(os.Args[1:]) {
	case 0:
		runPrompt()
	case 1:
		runFile(os.Args[1])
	default:
		fmt.Println("Usage: lox [script]")
		os.Exit(64)
	}
}

func runFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(74)
	}
	run(string(bytes))
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	fmt.Println("Welcome to Lox! Press Ctrl+C to exit.")
	for {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		run(line)
		hadError = false
	}
}

func run(source string) {
	scanner := lox.NewScanner(source)
	tokens := scanner.ScanTokens()
	parser := lox.NewParser(tokens)
	expression, err := parser.Parse()
	if err != nil {
		hadError = true
		return
	}

	fmt.Println(expression.AstPrint())
}
