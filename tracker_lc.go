package lexy

type LineColumnTracker struct {
	line   int
	column int
}

func NewLineColumnTracker() *LineColumnTracker {
	return &LineColumnTracker{}
}

func (l *LineColumnTracker) OnScan(r rune) {
	if r == '\n' {
		l.column = 0
		l.line++
	} else {
		l.column++
	}
}

func (l *LineColumnTracker) OnEmit(t *Token) {
	t.Data["line"] = l.line
	t.Data["column"] = l.column
}
