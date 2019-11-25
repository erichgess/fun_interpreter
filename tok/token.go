package tok

import (
	"fmt"
	"unicode"
)

type tokenType int

const (
	intType          tokenType = iota
	operatorType     tokenType = iota
	lParen           tokenType = iota
	rParen           tokenType = iota
	labelType        tokenType = iota
	assignmentOpType tokenType = iota
	commaType        tokenType = iota
)

type token struct {
	value string
	ty    tokenType
}

type used struct{}

type tokenizer struct {
	operatorRuneSet map[rune]used
	operators       []string
}

func newTokenizer(operators []string) tokenizer {
	tokenizer := tokenizer{
		operatorRuneSet: make(map[rune]used),
		operators:       operators,
	}

	// create operator rune set
	for _, op := range operators {
		for _, c := range []rune(op) {
			tokenizer.operatorRuneSet[c] = used{}
		}
	}

	return tokenizer
}

func (t *tokenizer) tokenize(text string) ([]token, error) {
	raw := []rune(text)
	tokens := make([]token, 0)
	// while not EOL
	for currentChar := 0; currentChar < len(raw); {
		// create a new token
		var tok token
		var err error
		tok, currentChar, err = t.extractToken(raw, currentChar)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
	}

	return tokens, nil
}

func (t *tokenizer) extractToken(raw []rune, currentChar int) (tok token, charPos int, err error) {
	// consume any whitespace
	for ; currentChar < len(raw) && unicode.IsSpace(raw[currentChar]); currentChar++ {
	}

	// Check the current char to determine what type of token this is
	// if char is digit then extract integer token
	if unicode.IsDigit(raw[currentChar]) {
		return t.extractIntToken(raw, currentChar)
	} else if unicode.IsLetter(raw[currentChar]) {
		return t.extractLabelToken(raw, currentChar)
	} else if _, ok := t.operatorRuneSet[raw[currentChar]]; ok {
		// if char is not then consume operator
		return t.extractOperatorToken(raw, currentChar)
	} else if raw[currentChar] == '(' {
		return token{
			value: "(",
			ty:    lParen,
		}, currentChar + 1, nil
	} else if raw[currentChar] == ')' {
		return token{
			value: ")",
			ty:    rParen,
		}, currentChar + 1, nil
	} else if raw[currentChar] == '=' {
		return token{
			value: "=",
			ty:    assignmentOpType,
		}, currentChar + 1, nil
	} else if raw[currentChar] == ',' {
		return token{
			value: ",",
			ty:    commaType,
		}, currentChar + 1, nil
	} else {
		return token{}, -1, fmt.Errorf("unexpected character during tokenization: %s", string(raw[currentChar]))
	}
}

func (t *tokenizer) extractIntToken(raw []rune, currentChar int) (tok token, charPos int, err error) {
	for charPos = currentChar; charPos < len(raw) && unicode.IsDigit(raw[charPos]); charPos++ {
	}

	tok = token{
		value: string(raw[currentChar:charPos]),
		ty:    intType,
	}

	return tok, charPos, nil
}

func (t *tokenizer) extractOperatorToken(raw []rune, currentChar int) (tok token, newCharPos int, err error) {
	charPos := currentChar
	for ; charPos < len(raw); charPos++ {
		if _, ok := t.operatorRuneSet[raw[charPos]]; !ok {
			break
		}
	}

	tok = token{
		value: string(raw[currentChar:charPos]),
		ty:    operatorType,
	}

	return tok, charPos, nil
}

func (t *tokenizer) extractLabelToken(raw []rune, currentChar int) (tok token, newCharPos int, err error) {
	if !unicode.IsLetter(raw[currentChar]) {
		return token{}, currentChar, fmt.Errorf("`label` must start with letter")
	}

	start := currentChar
	for ; currentChar < len(raw); currentChar++ {
		if !unicode.IsLetter(raw[currentChar]) && !unicode.IsDigit(raw[currentChar]) {
			break
		}
	}

	tok = token{
		value: string(raw[start:currentChar]),
		ty:    labelType,
	}

	return tok, currentChar, nil
}
