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

	interpreter.Execute("8 - 2 * 3")
	set := "test = 5 * 2"
	interpreter.Execute(set)
	input := "second = -3 * 4 - 2*test"
	interpreter.Execute(input)
	input2 := "second + 10"
	r, _ := interpreter.Execute(input2)
	fmt.Printf("second + 10 = %d\n", r)

	f := "def f x y = y * x"
	interpreter.Execute(f)
	g := "def g x = x + 2"
	interpreter.Execute(g)
	f = "f(6/2, g(3)) * 3"
	fmt.Println(interpreter.Execute(f))
}
