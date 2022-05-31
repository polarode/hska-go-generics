package testable

// Counting words in a string.
// A word consists of a sequence of bytes separated by blanks (' ').

func skip(p func(byte) bool) func(string) string {
	return func(s string) string {
		switch {
		case len(s) == 0:
			return ""
		case p(s[0]):
			return skip(p)(s[1:])
		default:
			return s
		}
	}
}
func Count(s string) int {
	skipBlanks := skip(func(b byte) bool {
		return b == ' '
	})

	skipWord := skip(func(b byte) bool {
		return b != ' '
	})
	switch {
	case len(s) == 0:
		return 0
	case s[0] == ' ':
		return Count(skipBlanks(s))
	default:
		return 1 + Count(skipWord(s))
	}
}
