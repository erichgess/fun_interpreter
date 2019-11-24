package tok

import "unicode"

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

func (t *tokenizer) tokenize(text string) []token {
	raw := []rune(text)
	tokens := make([]token, 0)
	// while not EOL
	for currentChar := 0; currentChar < len(raw); {
		// create a new token
		var tok token
		tok, currentChar = t.extractToken(raw, currentChar)
		tokens = append(tokens, tok)
	}

	return tokens
}

func (t *tokenizer) extractToken(raw []rune, currentChar int) (tok token, charPos int) {
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
		}, currentChar + 1
	} else if raw[currentChar] == ')' {
		return token{
			value: ")",
			ty:    rParen,
		}, currentChar + 1
	} else if raw[currentChar] == '=' {
		return token{
			value: "=",
			ty:    assignmentOpType,
		}, currentChar + 1
	} else if raw[currentChar] == ',' {
		return token{
			value: ",",
			ty:    commaType,
		}, currentChar + 1
	} else {
		panic("unexpected character during tokenization: " + string(raw[currentChar]))
	}
}

func (t *tokenizer) extractIntToken(raw []rune, currentChar int) (tok token, charPos int) {
	for charPos = currentChar; charPos < len(raw) && unicode.IsDigit(raw[charPos]); charPos++ {
	}

	tok = token{
		value: string(raw[currentChar:charPos]),
		ty:    intType,
	}

	return tok, charPos
}

func (t *tokenizer) extractOperatorToken(raw []rune, currentChar int) (tok token, newCharPos int) {
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

	return tok, charPos
}

func (t *tokenizer) extractLabelToken(raw []rune, currentChar int) (tok token, newCharPos int) {
	if !unicode.IsLetter(raw[currentChar]) {
		panic("`label` must start with letter")
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

	return tok, currentChar
}
