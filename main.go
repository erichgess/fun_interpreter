package main

import (
	"erichgess/parser/tok"
	"fmt"
)

func main() {
	input := "3 * (1 + 2)"
	tokenizer := tok.NewTokenizer([]string{"*", "+"})
	tokens := tokenizer.Tokenize(input)
	fmt.Printf("%+v\n", tokens)
	interpreter := tok.NewInterpreter()
	interpreter.AddExpressionOp("+", func(a, b int) int { return a + b })
	interpreter.AddFactorOp("*", func(a, b int) int { return a * b })
	fmt.Println(interpreter.Expression(tokens, 0))
}
