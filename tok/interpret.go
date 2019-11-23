package tok

import "strconv"

/*
BNF
Statement := Assignment | Expression
Assignment := Label AssignOp Expression
Expression := Factor[ExpOp Expression]
Factor := Term [FactorOp Factor]
Term := Integer | Label | UnaryOp Term | LParen Expression RParen
Integer := Digit+
Label := Alpha[Alpha|Digit]+
*/

// BinaryOperator is a function which takes two integers and returns one
type BinaryOperator func(a, b int) int

// UnaryOperator is a function which takes one integer and returns one
type UnaryOperator func(a int) int

// Interpreter allows a user to define a set of operators that fit within the following grammar
// and then use the interpreter to compute the results of programs written in that language.
//
// Grammar:
//
// - Expression := Factor[ExpOp Expression]
//
// - Factor := Term [FactorOp Factor]
//
// - Term := Integer | UnaryOp Term | LParen Expression RParen
//
// - Integer := Digit+
//
type Interpreter struct {
	expOps        map[string]BinaryOperator
	factorOps     map[string]BinaryOperator
	unaryOps      map[string]UnaryOperator
	labelBindings map[string]int
}

// NewInterpreter configures a new Interpreter object and returns it
func NewInterpreter() Interpreter {
	return Interpreter{
		expOps:        make(map[string]BinaryOperator),
		factorOps:     make(map[string]BinaryOperator),
		unaryOps:      make(map[string]UnaryOperator),
		labelBindings: make(map[string]int),
	}
}

// AddExpressionOp will add a new expression level operator for the interpreter to use.
// If an operator already exists at the expression level with the given symbol, it will
// be replaced.  If an operator with this symbol exists in the Factor operator set
// then this will fail.
func (i *Interpreter) AddExpressionOp(symbol string, apply BinaryOperator) {
	// make sure the operator does not exist in the Factor set
	if _, ok := i.factorOps[symbol]; ok {
		panic("attempting to add operator to expression set when it is already in factor set")
	}
	i.expOps[symbol] = apply
}

// AddFactorOp will add a new operator with the given symbol to the Factor level
// of interpretation.  If an operator already exists with this symbol for Factor level
// it will be replaced.  If an operator with this symbol exists in the Expression operator
// set then this will fail.
func (i *Interpreter) AddFactorOp(symbol string, apply BinaryOperator) {
	// make sure the operator does not exist in the Expression set
	if _, ok := i.expOps[symbol]; ok {
		panic("attempting to add operator to factor set when it is already in expression set")
	}
	i.factorOps[symbol] = apply
}

// AddUnaryOp will add a unary operator that will be applied at the Term level
// of the language.
func (i *Interpreter) AddUnaryOp(symbol string, apply UnaryOperator) {
	i.unaryOps[symbol] = apply
}

// Execute will take a program that uses the interpreters defined language
// and attempt to compute it's result
func (i *Interpreter) Execute(text string) int {
	// construct a tokenizer
	tokenizer := i.createTokenizer()

	tokens := tokenizer.tokenize(text)

	var pos int
	var result int
	if len(tokens) >= 3 && tokens[0].ty == labelType && tokens[1].ty == assignmentOpType {
		result, pos = i.assignment(tokens, 0)
	} else {
		result, pos = i.expression(tokens, 0)
	}
	if pos != len(tokens) {
		panic("unexpected tokens in expression")
	}
	return result
}

func (i *Interpreter) createTokenizer() tokenizer {
	// build a list of the operators in this interpreter
	opsList := make([]string, 0, len(i.expOps)+len(i.factorOps))
	for k := range i.expOps {
		opsList = append(opsList, k)
	}
	for k := range i.factorOps {
		opsList = append(opsList, k)
	}
	// create tokenizer
	return newTokenizer(opsList)
}

func (i *Interpreter) assignment(tokens []token, currentPos int) (result int, pos int) {
	if tokens[currentPos].ty != labelType {
		panic("invalid left side in assignment")
	}

	label := tokens[currentPos].value
	currentPos++

	if tokens[currentPos].ty != assignmentOpType {
		panic("expecting assignment operator")
	}
	currentPos++
	result, pos = i.expression(tokens, currentPos)
	i.labelBindings[label] = result

	return result, pos
}

func (i *Interpreter) expression(tokens []token, currentPos int) (result int, pos int) {
	result, pos = i.factor(tokens, currentPos)

	if pos < len(tokens) && tokens[pos].ty == operatorType {
		if op, ok := i.expOps[tokens[pos].value]; ok {
			pos++
			r, p := i.expression(tokens, pos)
			result = op(result, r)
			pos = p
		}
	}

	if pos < len(tokens) && tokens[pos].ty != rParen {
		panic("unexpected token in expression")
	}

	return result, pos
}

func (i *Interpreter) factor(tokens []token, currentPos int) (result int, pos int) {
	result, currentPos = i.term(tokens, currentPos)

	if currentPos < len(tokens) {
		if tokens[currentPos].ty == operatorType {
			if op, ok := i.factorOps[tokens[currentPos].value]; ok {
				currentPos++
				r, p := i.factor(tokens, currentPos)
				result = op(result, r)
				currentPos = p
			}
		}
	}

	return result, currentPos
}

func (i *Interpreter) term(tokens []token, currentPos int) (result int, pos int) {
	if tokens[currentPos].ty == lParen {
		currentPos++
		result, currentPos = i.expression(tokens, currentPos)

		// consume right paren
		if currentPos >= len(tokens) || tokens[currentPos].ty != rParen {
			panic("expected right paren")
		}
		currentPos++
	} else if tokens[currentPos].ty == operatorType {
		// if the operator is not unary then something is wrong
		if op, ok := i.unaryOps[tokens[currentPos].value]; ok {
			currentPos++
			result, currentPos = i.term(tokens, currentPos)
			result = op(result)
		} else {
			panic("unexpected token in factor")
		}
	} else if tokens[currentPos].ty == intType {
		result, _ = strconv.Atoi(tokens[currentPos].value)
		currentPos++
	} else if tokens[currentPos].ty == labelType {
		result, currentPos = i.lookupLabel(tokens, currentPos)
	}

	return result, currentPos
}

func (i *Interpreter) lookupLabel(tokens []token, currentPos int) (result int, pos int) {
	if tokens[currentPos].ty != labelType {
		panic("attempting to look up label binding for token that is not a label")
	}

	label := tokens[currentPos].value
	if v, ok := i.labelBindings[label]; ok {
		return v, currentPos + 1
	}

	panic("could not find value for label: " + label)
}
