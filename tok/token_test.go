package tok

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WhiteSpaceAtEnd(t *testing.T) {
	text := "2  "
	tokenizer := newTokenizer([]string{})
	tokens, err := tokenizer.tokenize(text)
	assert.NoError(t, err)
	assert.Equal(t, token{value: "2", ty: intType}, tokens[0])
}

func Test_WhiteSpaceAtStart(t *testing.T) {
	text := "  2"
	tokenizer := newTokenizer([]string{})
	tokens, err := tokenizer.tokenize(text)
	assert.NoError(t, err)
	assert.Equal(t, token{value: "2", ty: intType}, tokens[0])
}

func Test_WhiteSpaceBothEnds(t *testing.T) {
	text := "  2  "
	tokenizer := newTokenizer([]string{})
	tokens, err := tokenizer.tokenize(text)
	assert.NoError(t, err)
	assert.Equal(t, token{value: "2", ty: intType}, tokens[0])
}
