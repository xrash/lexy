package lexy

import (
	"fmt"
)

type DebugTracker struct {
	
}

func NewDebugTracker() *DebugTracker {
	return &DebugTracker{}
}

func (l *DebugTracker) OnScan(r rune) {
	fmt.Println(string(r))
}

func (l *DebugTracker) OnEmit(t *Token) {
	
}
