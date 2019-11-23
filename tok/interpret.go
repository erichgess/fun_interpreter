package tok

import "strconv"

/*
BNF
Expression := Factor[PLUS Expression]
Factor := Integer [MULTIPLY Factor]|'(' Expression ')'
Integer := Digit+
*/

type BinaryOperator interface {
	Apply(a, b int) int
}

type Interpreter struct {
}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i *Interpreter) Expression(tokens []Token, currentPos int) (result int, pos int) {
	result, pos = i.Factor(tokens, currentPos)

	if pos < len(tokens) {
		if tokens[pos].ty == OperatorType && tokens[pos].value == "+" {
			pos++
			r, p := i.Expression(tokens, pos)
			result += r
			pos = p
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
		if tokens[currentPos].ty == OperatorType && tokens[currentPos].value == "*" {
			currentPos++
			r, p := i.Factor(tokens, currentPos)
			result *= r
			currentPos = p
		}
	}

	return result, currentPos
}
