package main

import (
	"erichgess/parser/tok"
	"fmt"
)

func main() {
	interpreter := tok.NewInterpreter()
	interpreter.AddExpressionOp("+", func(a, b int) int { return a + b })
	interpreter.AddExpressionOp("-", func(a, b int) int { return a - b })
	interpreter.AddFactorOp("*", func(a, b int) int { return a * b })
	interpreter.AddFactorOp("/", func(a, b int) int { return a / b })
	interpreter.AddUnaryOp("-", func(a int) int { return -a })
	interpreter.AddUnaryOp("--", func(a int) int { return a - 1 })

	set := "x = 5 * 2"
	interpreter.Execute(set)

	g := "def g = x + 2"
	fmt.Println(interpreter.Execute(g))
	r := "g()"
	fmt.Println(interpreter.Execute(r))
}
