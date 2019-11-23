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
	fmt.Println(tok.Expression(tokens, 0))
}
