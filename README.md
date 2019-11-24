# fun_interpreter
A fun and simple script interpreter written in Go.  A user can define functions, variables, and write epxressions which evalulate basic arithmetic and function calls.

# Interpreter Usage
The interpreter evaluates single statements at a time.  A statement can define a function, assign a value to a variable, or evaulate 
an expression.  Operators must also be defined by the user.

## Defining Operators
Operators are defined by binding a unary or binary function to a symbol.  If the operator is binary its precedence is
set by adding it to either the Factor or Expression level of evaluation.

Factor has a higher precedence than Expression.

### Examples
```
	interpreter.AddExpressionOp("+", func(a, b int) int { return a + b })
	interpreter.AddExpressionOp("-", func(a, b int) int { return a - b })
	interpreter.AddFactorOp("*", func(a, b int) int { return a * b })
	interpreter.AddFactorOp("/", func(a, b int) int { return a / b })
	interpreter.AddUnaryOp("-", func(a int) int { return -a })
	interpreter.AddUnaryOp("--", func(a int) int { return a - 1 })
```

## Evaluating an Expression
Once operators are defined, expressions which use those operators can be evaulated:

```
	interpreter.Execute("8 - 2 * 3")
```

This will return the integer result of evaluation this expression.

## Assigning values to variables
A label can be assigned a value by using the assignment operator

```
	set := "test = 5 * 2"
	interpreter.Execute(set)
	input := "second = -3 * 4 - 2*test"
	interpreter.Execute(input)
	input2 := "second + 10"
	fmt.Println(interpreter.Execute(input2))
```

This will print the result `-22` to stdout.

## Defining Functions
Functions can also be defined by using the `def` symbol followed by the name of the function and, optionally, a list
of parameter labels.

```
	f := "def f x y = y * x"
	interpreter.Execute(f)
	g := "def g x = x + 2"
	interpreter.Execute(g)
	f = "f(6/2, g(3)) * 3"
	fmt.Println(interpreter.Execute(f))
```

This will print `45` to stdout.
