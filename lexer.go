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

type Lexer struct {
	tokens   chan *Token
	trackers []Tracker

	// metadata
	bag []rune
}

func NewLexer(tokens chan *Token) *Lexer {
	return &Lexer{
		tokens:   tokens,
		trackers: make([]Tracker, 0),
		bag:      make([]rune, 0),
	}
}

func (l *Lexer) AddTracker(t Tracker) {
	l.trackers = append(l.trackers, t)
}

func (l *Lexer) Collect(r ...rune) {
	if len(r) > 0 {
		l.bag = append(l.bag, r...)
	}
}

func (l *Lexer) Emit(key string) {
	t := newToken(key, string(l.bag))

	for _, tracker := range l.trackers {
		tracker.OnEmit(t)
	}

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

		for _, tracker := range l.trackers {
			tracker.OnScan(r)
		}

		if r == utf8.RuneError {
			msg := fmt.Sprintf("Error decoding rune %v (width %d)", r, width)
			l.tokens <- newToken(ErrorToken, msg)
			return fmt.Errorf(msg)
		}

		state, err = state(l, r)
		if err != nil {
			msg := fmt.Sprintf("Error: %v", err)
			l.tokens <- newToken(ErrorToken, msg)
			return fmt.Errorf(msg)
		}
	}

	state, err = state(l, 0)
	if err != nil {
		msg := fmt.Sprintf("Error: %v", err)
		l.tokens <- newToken(ErrorToken, msg)
		return fmt.Errorf(msg)
	}

	l.tokens <- newToken(EOFToken, "")

	return nil
}

func (l *Lexer) DoString(input string, state State) error {
	return l.Do(strings.NewReader(input), state)
}
