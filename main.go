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

	set := "test = 5 * 2"
	fmt.Println(interpreter.Execute(set))
	input := "second = -3 * 4 - 2*test"
	fmt.Println(interpreter.Execute(input))
	input2 := "second + 10"
	fmt.Println(interpreter.Execute(input2))

	f := "def f x y = y * x"
	interpreter.Execute(f)
	g := "def g x = x + 2"
	interpreter.Execute(g)
	f = "f(6/2, g(3)) * 3"
	fmt.Println(interpreter.Execute(f))
}
