package lexical

import (
	"strings"
	"unicode"
)

func Tokenize(s string) []string {
	fields := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	var terms []string
	for _, f := range fields {
		if len(f) < 2 {
			continue
		}
		terms = append(terms, f)
	}
	return terms
}
