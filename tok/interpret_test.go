package tok

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExpOperator(t *testing.T) {
	interpreter := NewInterpreter()
	interpreter.AddExpressionOp("+", func(a, b int) int { return a + b })
	input := "2 + 2"
	result, err := interpreter.Execute(input)

	assert.NoError(t, err)
	assert.Equal(t, 4, result)
}

func Test_FactorOperator(t *testing.T) {
	interpreter := NewInterpreter()
	interpreter.AddFactorOp("*", func(a, b int) int { return a * b })
	input := "2 * 2"
	result, err := interpreter.Execute(input)

	assert.NoError(t, err)
	assert.Equal(t, 4, result)
}

func Test_UnaryOperator(t *testing.T) {
	interpreter := NewInterpreter()
	interpreter.AddUnaryOp("-", func(a int) int { return -a })
	input := "-2"
	result, err := interpreter.Execute(input)

	assert.NoError(t, err)
	assert.Equal(t, -2, result)
}

func Test_OrderOperations(t *testing.T) {
	interpreter := NewInterpreter()
	interpreter.AddExpressionOp("+", func(a, b int) int { return a + b })
	interpreter.AddFactorOp("*", func(a, b int) int { return a * b })
	interpreter.AddUnaryOp("-", func(a int) int { return -a })
	input := "-2 * 3 + 1*(2*3)"
	result, _ := interpreter.Execute(input)

	assert.Equal(t, 0, result)
}

func Test_AddExpressionOpWhenFactorOpAlreadyExists_IsError(t *testing.T) {
	interpreter := NewInterpreter()
	interpreter.AddFactorOp("*", func(a, b int) int { return a * b })
	err := interpreter.AddExpressionOp("*", func(a, b int) int { return a + b })
	assert.Error(t, err)
}

func Test_AddFactorOpWhenExpressionOpAlreadyExists_IsError(t *testing.T) {
	interpreter := NewInterpreter()
	interpreter.AddExpressionOp("*", func(a, b int) int { return a + b })
	err := interpreter.AddFactorOp("*", func(a, b int) int { return a * b })
	assert.Error(t, err)
}

func Test_EvaluationExpressionWithUndefinedBinaryOp_IsError(t *testing.T) {
	interpreter := NewInterpreter()
	_, err := interpreter.Execute("2 + 2")
	assert.Error(t, err)
}

func Test_EvaluationExpressionWithUndefinedUnaryOp_IsError(t *testing.T) {
	interpreter := NewInterpreter()
	_, err := interpreter.Execute("-2")
	assert.Error(t, err)
}

func Test_EvaluateExpressionWithBinaryOpMissingRightOperand_IsError(t *testing.T) {
	interpreter := NewInterpreter()
	interpreter.AddExpressionOp("+", func(a, b int) int { return a + b })
	_, err := interpreter.Execute("2 + ")
	assert.Error(t, err)
}

func Test_EvaluateExpressionWithBinaryOpMissingLeftOperand_IsError(t *testing.T) {
	interpreter := NewInterpreter()
	interpreter.AddExpressionOp("+", func(a, b int) int { return a + b })
	_, err := interpreter.Execute("+ 2")
	assert.Error(t, err)
}

func Test_RightParenWithNoMatchingLeftParen_IsError(t *testing.T) {
	interpreter := NewInterpreter()
	_, err := interpreter.Execute("5)")
	assert.Error(t, err)
}

func Test_LeftParenWithNoMatchingRightParen_IsError(t *testing.T) {
	interpreter := NewInterpreter()
	_, err := interpreter.Execute("(5")
	assert.Error(t, err)
}

func Test_AssignValue(t *testing.T) {
	i := NewInterpreter()
	v, err := i.Execute("x = 5")
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
}

func Test_UseVariable(t *testing.T) {
	i := NewInterpreter()
	i.AddFactorOp("*", func(a, b int) int { return a * b })
	i.Execute("x = 5")
	v, err := i.Execute("2 * x")
	assert.NoError(t, err)
	assert.Equal(t, 10, v)
}

func Test_UseUndefinedVariable_IsError(t *testing.T) {
	i := NewInterpreter()
	i.AddFactorOp("*", func(a, b int) int { return a * b })
	_, err := i.Execute("2 * x")
	assert.Error(t, err)
}

func Test_OperatorIsAnOperand_IsError(t *testing.T) {
	i := NewInterpreter()
	i.AddFactorOp("*", func(a, b int) int { return a * b })
	_, err := i.Execute("2 * *")
	assert.Error(t, err)
}

func Test_TwoIntsInARow_IsError(t *testing.T) {
	i := NewInterpreter()
	_, err := i.Execute("2 2")
	assert.Error(t, err)
}

func Test_TwoVarsInARow_IsError(t *testing.T) {
	i := NewInterpreter()
	i.Execute("x = 5")
	_, err := i.Execute("x 2")
	assert.Error(t, err)
}

func Test_DefineFunctionNoParameters(t *testing.T) {
	i := NewInterpreter()
	_, err := i.Execute("def f = 5")
	assert.NoError(t, err)
}

func Test_DefineFunctionNoParametersUndefinedVar_NoError(t *testing.T) {
	i := NewInterpreter()
	_, err := i.Execute("def f = x")
	assert.Error(t, err)
}

func Test_DefineFunctionWithParameters(t *testing.T) {
	i := NewInterpreter()
	_, err := i.Execute("def f x = x")
	assert.NoError(t, err)
}

func Test_DefineFunctionWithInvalidParameterLabels_Error(t *testing.T) {
	i := NewInterpreter()
	for _, input := range []string{
		"def f 5 = x",
		"def f * = x",
		"def f -b = x",
		"def f b- = x",
		"def f ( = x",
		"def f ) = x",
	} {
		_, err := i.Execute(input)
		assert.Error(t, err)
	}
}

func Test_CallFunctionNoParameters(t *testing.T) {
	i := NewInterpreter()
	_, err := i.Execute("def f = 5")

	r, err := i.Execute("f()")
	assert.NoError(t, err)
	assert.Equal(t, 5, r)
}

func Test_CallFunctionWithParameters(t *testing.T) {
	i := NewInterpreter()
	i.AddFactorOp("*", func(a, b int) int { return a * b })
	_, err := i.Execute("def f x = x*2")
	assert.NoError(t, err)

	r, err := i.Execute("f(3)")
	assert.NoError(t, err)
	assert.Equal(t, 3, r)
}

func Test_CallFunctionMissingParameters_IsError(t *testing.T) {
	i := NewInterpreter()
	i.AddFactorOp("*", func(a, b int) int { return a * b })
	_, err := i.Execute("def f x = x*2")
	assert.NoError(t, err)

	_, err = i.Execute("f()")
	assert.Error(t, err)
}
