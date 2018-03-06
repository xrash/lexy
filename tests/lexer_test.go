package lexy

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/xrash/lexy"
)

var __tcase1 = `
1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 {This opening is called the Ruy Lopez.}
4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7
11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 15. Nb1 h6 16. Bh4 c5 17. dxe5
Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 Nxd6 21. Nc4 Nxc4 22. Bxc4 Nb6
23. Ne5 Rae8 24. Bxf7+ Rxf7 25. Nxf7 Rxe1+ 26. Qxe1 Kxf7 27. Qe3 Qg5 28. Qxg5
hxg5 29. b3 Ke6 30. a3 Kd6 31. axb4 cxb4 32. Ra5 Nd5 33. f3 Bc8 34. Kf2 Bf5
35. Ra7 g6 36. Ra6+ Kc5 37. Ke1 Nf4 38. g3 Nxh3 39. Kd2 Kb5 40. Rd6 Kc5 41. Ra6
Nf2 42. g4 Bd3 43. Re6 1/2-1/2
`

var __tcase2 = `
1. e4 { [%eval 0.31] } 1... e6 { [%eval 0.27] } 2. Nf3 { [%eval 0.13] } 2... d5 { [%eval 0.18] } 3. exd5 { [%eval 0.12] } 3... exd5 { [%eval 0.17] } 4. d4 { [%eval 0.08] } 4... Bd6 { [%eval 0.23] } 5. Bd3 { [%eval 0.07] } 5... Ne7 { [%eval 0.11] } 6. O-O { [%eval 0.1] } 6... Bf5 { [%eval 0.46] } 7. Bxf5 { [%eval 0.43] } 7... Nxf5 { [%eval 0.4] } 8. Qd3 { [%eval 0.46] } 8... Qd7 { [%eval 0.4] } 9. Re1+ { [%eval 0.5] } 9... Ne7 { [%eval 0.41] } 10. Ne5 { [%eval 0.07] } 10... Bxe5 { [%eval 0.16] } 11. dxe5 { [%eval 0.11] } 11... c6 { [%eval 0.16] } 12. c4 { [%eval 0.27] } 12... dxc4 { [%eval 0.37] } 13. Qxc4 { [%eval 0.43] } 13... Nd5 { [%eval 0.54] } 14. Nc3 { [%eval 0.54] } 14... Qe6 { [%eval 0.68] } 15. Ne4 { [%eval 0.36] } 15... Nd7? { [%eval 1.65] } 16. Nd6+ { [%eval 1.7] } 16... Ke7 { [%eval 2.15] } 17. Nxb7 { [%eval 1.71] } 17... Nxe5? { [%eval 3.69] } 18. Qc5+ { [%eval 4.22] } 18... Kd7 { [%eval 3.63] } 19. Bf4?! { [%eval 2.89] } 19... f6 { [%eval 3.05] } 20. Bxe5 { [%eval 2.56] } 20... fxe5 { [%eval 2.91] } 21. Rad1?! { [%eval 2.2] } 21... Rab8?? { [%eval 5.25] } 22. Qxa7 { [%eval 4.79] } 22... Ra8?? { [%eval 16.05] } 23. Nc5+ { [%eval 16.2] } 1-0
`

var __tcase3 = `
1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 {This opening is called the Ruy Lopez.}
4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7
11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 15. Nb1 h6 16. Bh4 c5 17. dxe5
Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 Nxd6 21. Nc4 Nxc4 22. Bxc4 Nb6
23. Ne5 Rae8 24. Bxf7+ Rxf7 25. Nxf7 Rxe1+ 26. Qxe1 Kxf7 27. Qe3 Qg5 28. Qxg5
hxg5 29. b3 Ke6 30. a3 Kd6 31. axb4 cxb4 32. Ra5 Nd5 33. f3 Bc8 34. Kf2 Bf5
35. Ra7 g6 36. Ra6+ Kc5 37. Ke1 Nf4 38. g3 Nxh3 39. Kd2 Kb5 40. Rd6 Kc5 41. Ra6
Nf2 42. g4 Bd3 43. Re6 1/2-1/2
`

func searchingNumberOrCommandOrResult(l *lexy.Lexer, r rune) (lexy.State, error) {
	if lexy.IsEOF(r) {
		return searchingNumberOrCommandOrResult, nil
	}

	if lexy.IsBlank(r) {
		return searchingNumberOrCommandOrResult, nil
	}

	if r == '{' {
		return inBraceCommentary, nil
	}

	if r == '*' {
		l.Collect(r)
		l.Emit("Result")
		return searchingNumberOrCommandOrResult, nil
	}

	if lexy.IsNumber(r) {
		l.Collect(r)
		return inNumberOrResult, nil
	}

	if lexy.IsLetter(r) {
		l.Collect(r)
		return inCommand, nil
	}

	return nil, fmt.Errorf("Expecting a number, a command or a result, got %v", r)
}

func inBraceCommentary(l *lexy.Lexer, r rune) (lexy.State, error) {
	if r == '}' {
		return searchingNumberOrCommandOrResult, nil
	}

	return inBraceCommentary, nil
}

func inNumberOrResult(l *lexy.Lexer, r rune) (lexy.State, error) {
	if lexy.IsBlank(r) {
		l.Emit("Number")
		return searchingNumberOrCommandOrResult, nil
	}

	if r == '/' || r == '-' {
		l.Collect(r)
		return inResult, nil
	}

	if lexy.IsNumber(r) || r == '.' {
		l.Collect(r)
		return inNumberOrResult, nil
	}

	return nil, fmt.Errorf("Expecting a number or a result, got %v", r)
}

func inResult(l *lexy.Lexer, r rune) (lexy.State, error) {
	if lexy.IsBlank(r) {
		l.Emit("Result")
		return searchingNumberOrCommandOrResult, nil
	}

	if lexy.IsNumber(r) || r == '-' || r == '/' {
		l.Collect(r)
		return inResult, nil
	}

	return nil, fmt.Errorf("Expecting a result, got %v", r)
}

func inCommand(l *lexy.Lexer, r rune) (lexy.State, error) {
	if lexy.IsBlank(r) {
		l.Emit("Command")
		return searchingNumberOrCommandOrResult, nil
	}

	l.Collect(r)

	return inCommand, nil
}

func TestDo(t *testing.T) {
	tokens := make(chan *lexy.Token, 1000)
	l := lexy.NewLexer(tokens)
	l.AddTracker(lexy.NewLineColumnTracker())
	err := l.DoString(__tcase3, searchingNumberOrCommandOrResult)
	assert.Nil(t, err)
}
