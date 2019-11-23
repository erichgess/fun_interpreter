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
	// make sure the operator does not exist in the Factor set
	if _, ok := i.factorOps[symbol]; ok {
		panic("attempting to add operator to expression set when it is already in factor set")
	}
	i.expOps[symbol] = apply
}

func (i *Interpreter) AddFactorOp(symbol string, apply BinaryOperator) {
	// make sure the operator does not exist in the Expression set
	if _, ok := i.expOps[symbol]; ok {
		panic("attempting to add operator to factor set when it is already in expression set")
	}
	i.factorOps[symbol] = apply
}

func (i *Interpreter) Execute(text string) int {
	// construct a tokenizer
	tokenizer := i.createTokenizer()

	tokens := tokenizer.Tokenize(text)

	result, pos := i.Expression(tokens, 0)
	if pos != len(tokens) {
		panic("unexpected tokens in expression")
	}
	return result
}

func (i *Interpreter) createTokenizer() Tokenizer {
	// build a list of the operators in this interpreter
	opsList := make([]string, 0, len(i.expOps)+len(i.factorOps))
	for k, _ := range i.expOps {
		opsList = append(opsList, k)
	}
	for k, _ := range i.factorOps {
		opsList = append(opsList, k)
	}
	// create tokenizer
	return NewTokenizer(opsList)
}

func (i *Interpreter) Expression(tokens []Token, currentPos int) (result int, pos int) {
	result, pos = i.Factor(tokens, currentPos)

	if pos < len(tokens) && tokens[pos].ty == OperatorType {
		if op, ok := i.expOps[tokens[pos].value]; ok {
			pos++
			r, p := i.Expression(tokens, pos)
			result = op(result, r)
			pos = p
		}
	}

	if pos < len(tokens) && tokens[pos].ty != RParen {
		panic("unexpected token in expression")
	}

	return result, pos
}

func (i *Interpreter) Factor(tokens []Token, currentPos int) (result int, pos int) {
	if tokens[currentPos].ty == LParen {
		currentPos++
		result, currentPos = i.Expression(tokens, currentPos)

		// consume right paren
		if currentPos >= len(tokens) || tokens[currentPos].ty != RParen {
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
