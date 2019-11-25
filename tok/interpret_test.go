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
	input := "-2 * 3 + 1"
	result, _ := interpreter.Execute(input)

	assert.Equal(t, -5, result)
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
