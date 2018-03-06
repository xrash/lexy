package lexy

type Tracker interface {
	OnScan(rune)
	OnEmit(*Token)
}
