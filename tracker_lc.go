package lexy

type LineColumnTracker struct {
	line       int
	column     int
	lastLineColumns int
}

func NewLineColumnTracker() *LineColumnTracker {
	return &LineColumnTracker{
		line: 1,
		column: 1,
	}
}

func (l *LineColumnTracker) OnScan(r rune) {
	if r == '\n' {
		l.line++
		l.lastLineColumns = l.column
		l.column = 1
	} else {
		l.column++
	}
}

func (l *LineColumnTracker) OnEmit(t *Token) {
	t.Data["line"] = l.line - 1
	t.Data["column"] = l.lastLineColumns
}
