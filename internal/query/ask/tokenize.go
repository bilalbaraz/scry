package ask

import "unicode"

// TokenizeQuery splits on non-alnum, lowercases, drops terms < 3 chars, and de-dupes.
func TokenizeQuery(q string) []string {
	seen := map[string]struct{}{}
	var tokens []string
	var buf []rune
	flush := func() {
		if len(buf) == 0 {
			return
		}
		term := string(buf)
		buf = buf[:0]
		if len(term) < 3 {
			return
		}
		if _, ok := seen[term]; ok {
			return
		}
		seen[term] = struct{}{}
		tokens = append(tokens, term)
	}

	for _, r := range q {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			buf = append(buf, unicode.ToLower(r))
			continue
		}
		flush()
	}
	flush()
	return tokens
}
