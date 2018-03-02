package lexy

type Token struct {
	Key   string
	Value string
	Data  map[string]interface{}
}

func newToken(key, value string) *Token {
	return &Token{
		Key:   key,
		Value: value,
		Data:  make(map[string]interface{}),
	}
}
