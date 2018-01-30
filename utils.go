package lexy

func IsSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func IsLineBreak(r rune) bool {
	return r == '\n'
}

func IsEOF(r rune) bool {
	return r == 0
}

func IsBlank(r rune) bool {
	return IsSpace(r) || IsLineBreak(r) || IsEOF(r)
}

func IsNumber(r rune) bool {
	return r >= 48 && r <= 57
}

func IsLowercaseLetter(r rune) bool {
	return r >= 97 && r <= 122
}

func IsUppercaseLetter(r rune) bool {
	return r >= 65 && r <= 90
}

func IsLetter(r rune) bool {
	return IsLowercaseLetter(r) || IsUppercaseLetter(r)
}
