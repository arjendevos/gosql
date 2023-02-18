package functions

func singularize(word string) string {
	if len(word) == 0 {
		return word
	}
	lastChar := word[len(word)-1]
	switch lastChar {
	case 's':
		if len(word) > 2 && word[len(word)-2] == 'e' && word[len(word)-3] == 'i' {
			return word[:len(word)-3] + "y"
		}
		return word[:len(word)-1]
	case 'x', 'z':
		return word[:len(word)-1]
	case 'h':
		if len(word) > 1 && (word[len(word)-2] == 's' || word[len(word)-2] == 'c') {
			return word[:len(word)-1]
		}
		return word
	case 'y':
		if len(word) > 1 && (word[len(word)-2] == 'a' || word[len(word)-2] == 'e' || word[len(word)-2] == 'i' || word[len(word)-2] == 'o' || word[len(word)-2] == 'u') {
			return word
		}
		return word[:len(word)-1] + "ies"
	default:
		return word
	}
}
