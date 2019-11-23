package main

import (
	"erichgess/parser/tok"
	"fmt"
)

func main() {
	input := "(3 * 4 - 2)/2"
	interpreter := tok.NewInterpreter()
	interpreter.AddExpressionOp("+", func(a, b int) int { return a + b })
	interpreter.AddExpressionOp("-", func(a, b int) int { return a - b })
	interpreter.AddFactorOp("*", func(a, b int) int { return a * b })
	interpreter.AddFactorOp("/", func(a, b int) int { return a / b })
	fmt.Println(interpreter.Execute(input))
}
