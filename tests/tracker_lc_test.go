package lexy

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/xrash/lexy"
)

func TestLineColumnTracker(t *testing.T) {

	// Input data.
	data := `
second line

fourth line
`

	// The only state.
	var searchingLineBreak func(l *lexy.Lexer, r rune) (lexy.State, error)
	searchingLineBreak = func(l *lexy.Lexer, r rune) (lexy.State, error) {
		l.Collect(r)

		if lexy.IsLineBreak(r) {
			l.Emit("Line")
		}

		return searchingLineBreak, nil
	}

	// Configure the lexer.
	tokens := make(chan *lexy.Token, 1000)
	l := lexy.NewLexer(tokens)
	l.AddTracker(lexy.NewLineColumnTracker())

	// Run stuff.
	err := l.DoString(data, searchingLineBreak)
	assert.Nil(t, err)

	firstLine := <-tokens
	assert.Equal(t, "\n", firstLine.Value)
	secondLine := <-tokens
	assert.Equal(t, "second line\n", secondLine.Value)
	thirdLine := <-tokens
	assert.Equal(t, "\n", thirdLine.Value)
	fourthLine := <-tokens
	assert.Equal(t, "fourth line\n", fourthLine.Value)
}
