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

func searchingAOrB(l *lexy.Lexer, r rune) (lexy.State, error) {
	if lexy.IsEOF(r) {
		return searchingAOrB, nil
	}

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
&{SequenceOfAs aaaa map[]}
&{SequenceOfBs bb map[]}
&{SequenceOfAs aaaaaaaa map[]}
&{SequenceOfBs bbbbbb map[]}
&{SequenceOfAs aaaa map[]}
&{EOF  map[]}
<nil>
```

This is the sequence of produced tokens, then the `nil` error.

## Trackers

A Tracker tracks all the lexer process with the functions `OnScan(rune)` and `OnEmit(*Token)`. The function `OnScan(rune)` runs for every scanned rune and the function `OnEmit(*Token)` runs for every token emitted.

```go
type Tracker interface {
	OnScan(rune)
	OnEmit(*Token)
}
```

The tracker can do whatever it wants, like count the columns and lines and then emit this information with the token using the `Token.Data` map. Trackers can also be used for debugging purposes. There are two Trackers already available in this package: the `LineColumnTracker` and the `DebugTracker`.

You can add trackers using the `AddTracker(Tracker)` method of the `Lexer`, for example:

```go
l := lexy.NewLexer(tokens)
l.AddTracker(lexy.NewLineColumnTracker())
```

Now, every token will have `Token.Data["line"]` and `Token.Data["column"]` available.
