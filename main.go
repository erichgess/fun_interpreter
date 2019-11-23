package main

import (
	"erichgess/parser/tok"
	"fmt"
)

func main() {
	input := "235*"
	tokenizer := tok.NewTokenizer([]string{"*", "+"})
	tokens := tokenizer.Tokenize(input)
	fmt.Printf("%+v\n", tokens)
}
