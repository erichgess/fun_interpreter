package tok

import "unicode"

type tokenType int

const (
	IntType      tokenType = iota
	OperatorType tokenType = iota
	LParen       tokenType = iota
	RParen       tokenType = iota
)

type Token struct {
	value string
	ty    tokenType
}

type used struct{}

type Tokenizer struct {
	operatorRuneSet map[rune]used
	operators       []string
}

func NewTokenizer(operators []string) Tokenizer {
	tokenizer := Tokenizer{
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

func (t *Tokenizer) Tokenize(text string) []Token {
	raw := []rune(text)
	tokens := make([]Token, 0)
	// while not EOL
	for currentChar := 0; currentChar < len(raw); {
		// create a new token
		var tok Token
		tok, currentChar = t.extractToken(raw, currentChar)
		tokens = append(tokens, tok)
	}

	return tokens
}

func (t *Tokenizer) extractToken(raw []rune, currentChar int) (token Token, charPos int) {
	// consume any whitespace
	for ; currentChar < len(raw) && unicode.IsSpace(raw[currentChar]); currentChar++ {
	}

	// Check the current char to determine what type of token this is
	// if char is digit then extract integer token
	if unicode.IsDigit(raw[currentChar]) {
		return t.extractIntToken(raw, currentChar)
	} else if _, ok := t.operatorRuneSet[raw[currentChar]]; ok {
		// if char is not then consume operator
		return t.extractOperatorToken(raw, currentChar)
	} else if raw[currentChar] == '(' {
		return Token{
			value: "(",
			ty:    LParen,
		}, currentChar + 1
	} else if raw[currentChar] == ')' {
		return Token{
			value: ")",
			ty:    RParen,
		}, currentChar + 1
	} else {
		panic("unexpected character during tokenization")
	}
}

func (t *Tokenizer) extractIntToken(raw []rune, currentChar int) (token Token, charPos int) {
	for charPos = currentChar; charPos < len(raw) && unicode.IsDigit(raw[charPos]); charPos++ {
	}

	tok := Token{
		value: string(raw[currentChar:charPos]),
		ty:    IntType,
	}

	return tok, charPos
}

func (t *Tokenizer) extractOperatorToken(raw []rune, currentChar int) (token Token, newCharPos int) {
	charPos := currentChar
	for ; charPos < len(raw); charPos++ {
		if _, ok := t.operatorRuneSet[raw[charPos]]; !ok {
			break
		}
	}

	tok := Token{
		value: string(raw[currentChar:charPos]),
		ty:    OperatorType,
	}

	return tok, charPos
}
