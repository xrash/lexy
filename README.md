# Lexy

Simple lexer for golang, below is an example of how to use it:

```go
package main

import (
	"fmt"
	"github.com/xrash/lexy"
)

var __data = `
aaaa bb aaaaaaaa bbbbbb aaaa
`

func consumeTokens(tokens chan *lexy.Token, back chan error) {
	for t := range tokens {
		fmt.Println(t)

		if t.Key == lexy.ErrorToken {
			back <- fmt.Errorf("Error while lexing: %v", t.Value)
			return
		}

		if t.Key == lexy.EOFToken {
			back <- nil
			return
		}
	}
}

func main() {
	tokens := make(chan *lexy.Token, 1000)
	back := make(chan error)

	go consumeTokens(tokens, back)

	l := lexy.NewLexer(tokens)
	err := l.DoString(__data, searchingAOrB)
	if err != nil {
		panic(err)
	}

	err = <-back

	fmt.Println(err)
}

func searchingAOrB(l *lexy.Lexer, r rune) (lexy.State, error) {
	if lexy.IsBlank(r) {
		return searchingAOrB, nil
	}

	if r == 'a' {
		l.Collect(r)
		return inA, nil
	}

	if r == 'b' {
		l.Collect(r)
		return inB, nil
	}

	return nil, fmt.Errorf("Expecting a or b, got %v", r)
}

func inA(l *lexy.Lexer, r rune) (lexy.State, error) {
	if r == 'a' {
		l.Collect(r)
		return inA, nil
	}

	if lexy.IsBlank(r) {
		l.Emit("SequenceOfAs")
		return searchingAOrB, nil
	}

	return nil, fmt.Errorf("Expectin a or end of sequence, got %v", r)
}

func inB(l *lexy.Lexer, r rune) (lexy.State, error) {
	if r == 'b' {
		l.Collect(r)
		return inB, nil
	}

	if lexy.IsBlank(r) {
		l.Emit("SequenceOfBs")
		return searchingAOrB, nil
	}

	return nil, fmt.Errorf("Expectin b or end of sequence, got %v", r)
}
```

This program should output:

```
&{SequenceOfAs aaaa 1 5}
&{SequenceOfBs bb 1 8}
&{SequenceOfAs aaaaaaaa 1 17}
&{SequenceOfBs bbbbbb 1 24}
&{SequenceOfAs aaaa 2 0}
&{EOF  -1 -1}
<nil>
```

This is the sequence of produced tokens, then the `nil` error.
