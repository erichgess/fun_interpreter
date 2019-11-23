package tok

import "strconv"

/*
BNF
Expression := Factor[PLUS Expression]
Factor := Integer [MULTIPLY Factor]
Integer := Digit+
*/

func Expression(tokens []Token, currentPos int) (result int, pos int) {
	result, pos = Factor(tokens, currentPos)

	if pos < len(tokens) {
		if tokens[pos].ty == OperatorType && tokens[pos].value == "+" {
			pos++
			r, p := Expression(tokens, pos)
			result += r
			pos = p
		} else {
			panic("unexpected token in expression")
		}
	}

	return result, pos
}

func Factor(tokens []Token, currentPos int) (result int, pos int) {
	if tokens[currentPos].ty != IntType {
		panic("unexpected token")
	}

	result, _ = strconv.Atoi(tokens[currentPos].value)
	pos = currentPos + 1

	if pos < len(tokens) {
		if tokens[pos].ty == OperatorType && tokens[pos].value == "*" {
			pos++
			r, p := Expression(tokens, pos)
			result *= r
			pos = p
		}
	}

	return result, pos
}
