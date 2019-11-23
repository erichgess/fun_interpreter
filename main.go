package main

import (
	"erichgess/parser/tok"
	"fmt"
)

func main() {
	input := "1 + 2 * 3"
	tokenizer := tok.NewTokenizer([]string{"*", "+"})
	tokens := tokenizer.Tokenize(input)
	fmt.Printf("%+v\n", tokens)
	fmt.Println(tok.Expression(tokens, 0))
}
