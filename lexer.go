package lexy

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

const (
	ErrorToken string = "ERROR"
	EOFToken   string = "EOF"
)

type State func(l *Lexer, r rune) (State, error)

type Token struct {
	Key    string
	Value  string
	Line   int
	Column int
}

func newToken(key, value string, line, column int) *Token {
	return &Token{
		Key:    key,
		Value:  value,
		Line:   line,
		Column: column,
	}
}

type Lexer struct {
	tokens chan *Token

	// metadata
	bag    []rune
	column int
	line   int
}

func NewLexer(tokens chan *Token) *Lexer {
	return &Lexer{
		tokens: tokens,
		bag:    make([]rune, 0),
	}
}

func (l *Lexer) Collect(r ...rune) {
	if len(r) > 0 {
		l.bag = append(l.bag, r...)
	}
}

func (l *Lexer) Emit(key string) {
	t := newToken(key, string(l.bag), l.line, l.column)
	l.tokens <- t
	l.bag = make([]rune, 0)
}

func (l *Lexer) Do(r io.Reader, state State) error {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanRunes)
	l.bag = make([]rune, 0)
	var err error

	for s.Scan() {
		r, width := utf8.DecodeRune(s.Bytes())

		if r == utf8.RuneError {
			msg := fmt.Sprintf("Error decoding rune %v (width %d)", r, width)
			l.tokens <- newToken(ErrorToken, msg, -1, -1)
			return fmt.Errorf(msg)
		}

		if r == '\n' {
			l.column = 0
			l.line++
		} else {
			l.column++
		}

		state, err = state(l, r)
		if err != nil {
			msg := fmt.Sprintf("Error: %v", err)
			l.tokens <- newToken(ErrorToken, msg, -1, -1)
			return fmt.Errorf(msg)
		}
	}

	state, err = state(l, 0)
	if err != nil {
		msg := fmt.Sprintf("Error: %v", err)
		l.tokens <- newToken(ErrorToken, msg, -1, -1)
		return fmt.Errorf(msg)
	}

	l.tokens <- newToken(EOFToken, "", -1, -1)

	return nil
}

func (l *Lexer) DoString(input string, state State) error {
	return l.Do(strings.NewReader(input), state)
}
