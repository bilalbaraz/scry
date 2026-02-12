package ask

import "strings"

// DistinctTermMatchCount counts distinct query terms present in text (case-insensitive).
func DistinctTermMatchCount(terms []string, text string) int {
	if len(terms) == 0 || text == "" {
		return 0
	}
	lower := strings.ToLower(text)
	count := 0
	seen := map[string]struct{}{}
	for _, t := range terms {
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		if strings.Contains(lower, t) {
			seen[t] = struct{}{}
			count++
		}
	}
	return count
}
