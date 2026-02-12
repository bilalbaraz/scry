package ask

import "strings"

// FilterByQueryTerms keeps only chunks whose text contains at least one query term.
func FilterByQueryTerms(chunks []Chunk, terms []string) []Chunk {
	if len(terms) == 0 || len(chunks) == 0 {
		return nil
	}
	var out []Chunk
	for _, ch := range chunks {
		text := strings.ToLower(ch.Text)
		if hasAnyTerm(text, terms) {
			out = append(out, ch)
		}
	}
	return out
}

func hasAnyTerm(text string, terms []string) bool {
	for _, t := range terms {
		if t == "" {
			continue
		}
		if strings.Contains(text, t) {
			return true
		}
	}
	return false
}
