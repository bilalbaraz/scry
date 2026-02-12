package ask

import "strings"

// SelectTopEvidence returns at most n chunks by score, stable by path/id.
func SelectTopEvidence(chunks []Chunk, n int) []Chunk {
	if n <= 0 || len(chunks) == 0 {
		return nil
	}
	if len(chunks) <= n {
		return chunks
	}
	return chunks[:n]
}

// SnippetAroundTerm returns a snippet around the first matched term, up to maxChars.
func SnippetAroundTerm(text string, terms []string, maxChars int) string {
	if maxChars <= 0 {
		return ""
	}
	lower := strings.ToLower(text)
	pos := -1
	for _, t := range terms {
		if t == "" {
			continue
		}
		idx := strings.Index(lower, t)
		if idx >= 0 {
			pos = idx
			break
		}
	}
	if pos < 0 {
		return trimSnippet(text, maxChars)
	}
	start := pos - maxChars/3
	if start < 0 {
		start = 0
	}
	end := start + maxChars
	if end > len(text) {
		end = len(text)
	}
	snippet := strings.TrimSpace(text[start:end])
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(text) {
		snippet = snippet + "..."
	}
	return snippet
}

func trimSnippet(text string, maxChars int) string {
	if len(text) <= maxChars {
		return strings.TrimSpace(text)
	}
	return strings.TrimSpace(text[:maxChars]) + "..."
}
