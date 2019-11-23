package tok

import "strconv"

/*
BNF
Expression := Factor[PLUS Expression]
Factor := Integer [MULTIPLY Factor]|'(' Expression ')'
Integer := Digit+
*/

type BinaryOperator func(a, b int) int

type Interpreter struct {
	expOps    map[string]BinaryOperator
	factorOps map[string]BinaryOperator
}

func NewInterpreter() Interpreter {
	return Interpreter{
		expOps:    make(map[string]BinaryOperator),
		factorOps: make(map[string]BinaryOperator),
	}
}

func (i *Interpreter) AddExpressionOp(symbol string, apply BinaryOperator) {
	i.expOps[symbol] = apply
}

func (i *Interpreter) AddFactorOp(symbol string, apply BinaryOperator) {
	i.factorOps[symbol] = apply
}

func (i *Interpreter) Expression(tokens []Token, currentPos int) (result int, pos int) {
	result, pos = i.Factor(tokens, currentPos)

	if pos < len(tokens) {
		if tokens[pos].ty == OperatorType {
			if op, ok := i.expOps[tokens[pos].value]; ok {
				pos++
				r, p := i.Expression(tokens, pos)
				result = op(result, r)
				pos = p
			}
		}
	}

	return result, pos
}

func (i *Interpreter) Factor(tokens []Token, currentPos int) (result int, pos int) {
	if tokens[currentPos].ty == LParen {
		currentPos++
		result, currentPos = i.Expression(tokens, currentPos)

		// consume right paren
		if tokens[currentPos].ty != RParen {
			panic("expected right paren")
		}
		currentPos++
	} else if tokens[currentPos].ty == IntType {
		result, _ = strconv.Atoi(tokens[currentPos].value)
		currentPos++
	}

	if currentPos < len(tokens) {
		if tokens[currentPos].ty == OperatorType {
			if op, ok := i.factorOps[tokens[currentPos].value]; ok {
				currentPos++
				r, p := i.Factor(tokens, currentPos)
				result = op(result, r)
				currentPos = p
			}
		}
	}

	return result, currentPos
}
