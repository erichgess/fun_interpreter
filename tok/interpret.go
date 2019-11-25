package tok

import (
	"fmt"
	"strconv"
)

/*
BNF
Statement := Assignment | Expression | FuncDef
FuncDef := Label(def) Label+ AssignOp Expression
Assignment := Label AssignOp Expression
Expression := Factor[ExpOp Expression]
Factor := Term [FactorOp Factor]
Term := Integer | Label | UnaryOp Term | LParen Expression RParen | Label LParen [Label[,Label]*] RParen
Integer := Digit+
Label := Alpha[Alpha|Digit]+
*/

// Interpreter allows a user to define a set of operators that fit within the following grammar
// and then use the interpreter to compute the results of programs written in that language.
//
// Grammar:
//
// - Expression := Factor[ExpOp Expression]
//
// - Factor := Term [FactorOp Factor]
//
// - Term := Integer | UnaryOp Term | LParen Expression RParen | Label LParen RParen
//
// - Integer := Digit+
//
type Interpreter struct {
	expOps        map[string]BinaryOperator
	factorOps     map[string]BinaryOperator
	unaryOps      map[string]UnaryOperator
	labelBindings map[string]int
	funcBindings  map[string]function
}

// BinaryOperator is a function which takes two integers and returns one
type BinaryOperator func(a, b int) int

// UnaryOperator is a function which takes one integer and returns one
type UnaryOperator func(a int) int

type function struct {
	tokens     []token
	parameters []string
	name       string
	expOps     map[string]BinaryOperator
	factorOps  map[string]BinaryOperator
	unaryOps   map[string]UnaryOperator
}

func (f *function) apply(params []int) (int, error) {
	interpreter := NewInterpreter()
	interpreter.expOps = f.expOps
	interpreter.factorOps = f.factorOps
	interpreter.unaryOps = f.unaryOps

	if len(params) != len(f.parameters) {
		return 0, fmt.Errorf("missing parameters; expected %d got %d", len(f.parameters), len(params))
	}

	// bind the parameter labels to their given values
	for i, label := range f.parameters {
		interpreter.labelBindings[label] = params[i]
	}

	return interpreter.executeTokens(f.tokens)
}

// NewInterpreter configures a new Interpreter object and returns it
func NewInterpreter() Interpreter {
	return Interpreter{
		expOps:        make(map[string]BinaryOperator),
		factorOps:     make(map[string]BinaryOperator),
		unaryOps:      make(map[string]UnaryOperator),
		labelBindings: make(map[string]int),
		funcBindings:  make(map[string]function),
	}
}

// AddExpressionOp will add a new expression level operator for the interpreter to use.
// If an operator already exists at the expression level with the given symbol, it will
// be replaced.  If an operator with this symbol exists in the Factor operator set
// then this will fail.
func (i *Interpreter) AddExpressionOp(symbol string, apply BinaryOperator) error {
	// make sure the operator does not exist in the Factor set
	if _, ok := i.factorOps[symbol]; ok {
		return fmt.Errorf("attempting to add operator to expression set when it is already in factor set")
	}
	i.expOps[symbol] = apply
	return nil
}

// AddFactorOp will add a new operator with the given symbol to the Factor level
// of interpretation.  If an operator already exists with this symbol for Factor level
// it will be replaced.  If an operator with this symbol exists in the Expression operator
// set then this will fail.
func (i *Interpreter) AddFactorOp(symbol string, apply BinaryOperator) error {
	// make sure the operator does not exist in the Expression set
	if _, ok := i.expOps[symbol]; ok {
		return fmt.Errorf("attempting to add operator to factor set when it is already in expression set")
	}
	i.factorOps[symbol] = apply
	return nil
}

// AddUnaryOp will add a unary operator that will be applied at the Term level
// of the language.
func (i *Interpreter) AddUnaryOp(symbol string, apply UnaryOperator) error {
	i.unaryOps[symbol] = apply
	return nil
}

// Execute will take a program that uses the interpreters defined language
// and attempt to compute it's result
func (i *Interpreter) Execute(text string) (int, error) {
	// construct a tokenizer
	tokenizer := i.createTokenizer()

	tokens, err := tokenizer.tokenize(text)

	if err != nil {
		return 0, err
	}

	return i.executeTokens(tokens)
}

func (i *Interpreter) executeTokens(tokens []token) (int, error) {
	var pos int
	var result int
	if len(tokens) >= 3 && tokens[0].ty == labelType && tokens[1].ty == assignmentOpType {
		result, pos = i.assignment(tokens, 0)
	} else if tokens[0].ty == labelType && tokens[0].value == "def" {
		var f function
		var err error
		f, pos, err = i.functionDef(tokens, 0)
		if err != nil {
			return 0, err
		}
		i.funcBindings[f.name] = f
	} else {
		var err error
		result, pos, err = i.expression(tokens, 0)
		if err != nil {
			return 0, err
		}
	}
	if pos != len(tokens) {
		return 0, fmt.Errorf("unexpected tokens in expression: %s", tokens[pos].value)
	}
	return result, nil
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
	for k := range i.unaryOps {
		opsList = append(opsList, k)
	}
	// create tokenizer
	return newTokenizer(opsList)
}

func (i *Interpreter) functionDef(tokens []token, currentPos int) (f function, pos int, err error) {
	if tokens[currentPos].ty != labelType || tokens[currentPos].value != "def" {
		panic("unexpected token")
	}
	currentPos++

	if tokens[currentPos].ty != labelType {
		panic("expected label token found " + tokens[currentPos].value)
	}
	funcName := tokens[currentPos].value
	currentPos++

	// each label from now until an assignment operator is encountered is a function parameter
	parameters := make([]string, 0)
	for ; currentPos < len(tokens) && tokens[currentPos].ty == labelType; currentPos++ {
		parameters = append(parameters, tokens[currentPos].value)
	}

	// consume assignment operator
	if tokens[currentPos].ty != assignmentOpType {
		return function{}, currentPos, fmt.Errorf("expected '=' found '%s'", tokens[currentPos].value)
	}
	currentPos++

	// the remaining tokens are the function logic
	funcTokens := tokens[currentPos:]

	err = checkFunctionCorrectness(parameters, funcTokens)
	if err != nil {
		return function{}, currentPos, err
	}

	return function{
		name:       funcName,
		tokens:     funcTokens,
		parameters: parameters,
		expOps:     i.expOps,
		factorOps:  i.factorOps,
		unaryOps:   i.unaryOps,
	}, len(tokens), err
}

func checkFunctionCorrectness(parameters []string, tokens []token) error {
	// convert parameters into look up table
	paramLookup := make(map[string]bool)
	for _, p := range parameters {
		paramLookup[p] = true
	}

	// check that any variable in the function definition has a corresponding parameter
	for _, tok := range tokens {
		if tok.ty == labelType {
			if _, ok := paramLookup[tok.value]; !ok {
				return fmt.Errorf("undefined variable: %s", tok.value)
			}
		}
	}

	return nil
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
	result, pos, _ = i.expression(tokens, currentPos)
	i.labelBindings[label] = result

	return result, pos
}

func (i *Interpreter) expression(tokens []token, currentPos int) (result int, pos int, err error) {
	result, pos, err = i.factor(tokens, currentPos)
	if err != nil {
		return result, pos, err
	}

	if pos < len(tokens) && tokens[pos].ty == operatorType {
		if op, ok := i.expOps[tokens[pos].value]; ok {
			pos++
			r, p, err := i.expression(tokens, pos)
			if err != nil {
				return r, p, err
			}
			result = op(result, r)
			pos = p
		}
	}

	if pos < len(tokens) && (tokens[pos].ty != rParen && tokens[pos].ty != commaType) {
		return result, pos, fmt.Errorf("unexpected token in expression: %s", tokens[pos].value)
	}

	return result, pos, nil
}

func (i *Interpreter) factor(tokens []token, currentPos int) (result int, pos int, err error) {
	result, currentPos, err = i.term(tokens, currentPos)
	if err != nil {
		return result, currentPos, err
	}

	if currentPos < len(tokens) {
		if tokens[currentPos].ty == operatorType {
			if op, ok := i.factorOps[tokens[currentPos].value]; ok {
				currentPos++
				r, p, err := i.factor(tokens, currentPos)
				if err != nil {
					return r, p, err
				}
				result = op(result, r)
				currentPos = p
			}
		}
	}

	return result, currentPos, nil
}

func (i *Interpreter) term(tokens []token, currentPos int) (result int, pos int, err error) {
	if currentPos == len(tokens) {
		return 0, currentPos, fmt.Errorf("expecting term, but none found")
	}
	if tokens[currentPos].ty == lParen {
		currentPos++
		result, currentPos, err = i.expression(tokens, currentPos)
		if err != nil {
			return result, currentPos, err
		}

		// consume right paren
		if currentPos >= len(tokens) || tokens[currentPos].ty != rParen {
			return 0, currentPos, fmt.Errorf("expected right paren")
		}
		currentPos++
	} else if tokens[currentPos].ty == operatorType {
		// if the operator is not unary then something is wrong
		if op, ok := i.unaryOps[tokens[currentPos].value]; ok {
			currentPos++
			result, currentPos, err = i.term(tokens, currentPos)
			if err != nil {
				return 0, currentPos, err
			}
			result = op(result)
		} else {
			return 0, 0, fmt.Errorf("unexpected token in factor: %s", tokens[currentPos].value)
		}
	} else if tokens[currentPos].ty == intType {
		result, _ = strconv.Atoi(tokens[currentPos].value)
		currentPos++
	} else if tokens[currentPos].ty == labelType {
		// check if this is a function call
		if len(tokens)-currentPos-1 >= 1 && tokens[currentPos+1].ty == lParen {
			result, currentPos, err = i.functionCall(tokens, currentPos)
		} else {
			result, currentPos, err = i.lookupLabel(tokens, currentPos)
		}
	}

	return result, currentPos, err
}

func (i *Interpreter) functionCall(tokens []token, currentPos int) (result int, pos int, err error) {
	funcName := tokens[currentPos].value
	currentPos++
	if tokens[currentPos].ty != lParen {
		return 0, 0, fmt.Errorf("expected lparen")
	}
	currentPos++

	// Get function parameters
	params := make([]int, 0)
	for currentPos < len(tokens) && tokens[currentPos].ty != rParen {
		var v int
		v, currentPos, _ = i.expression(tokens, currentPos)
		params = append(params, v)

		if tokens[currentPos].ty == commaType {
			currentPos++
		}
	}

	if tokens[currentPos].ty != rParen {
		return 0, 0, fmt.Errorf("expected rparen")
	}
	currentPos++
	if f, ok := i.funcBindings[funcName]; ok {
		result, err = f.apply(params)
		if err != nil {
			return 0, 0, err
		}
	} else {
		return 0, 0, fmt.Errorf("function name not found: " + funcName)
	}

	return result, currentPos, nil
}

func (i *Interpreter) lookupLabel(tokens []token, currentPos int) (result int, pos int, err error) {
	if tokens[currentPos].ty != labelType {
		panic("attempting to look up label binding for token that is not a label")
	}

	label := tokens[currentPos].value
	if v, ok := i.labelBindings[label]; ok {
		return v, currentPos + 1, nil
	}

	return 0, currentPos, fmt.Errorf("could not find value for label: " + label)
}
