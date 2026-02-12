package ask

import "strings"

type BoostRule struct {
	Terms []string
	Paths []string
	Bonus float64
}

func DefaultBoostRules() []BoostRule {
	return []BoostRule{
		{
			Terms: []string{"ignore", "scan", "gitignore", "pattern", "exclude"},
			Paths: []string{"pkg/ignore/", "pkg/scan/"},
			Bonus: 2.0,
		},
	}
}

// ApplyBoosts adds score bonuses based on query terms and file path.
func ApplyBoosts(chunks []Chunk, terms []string, rules []BoostRule) []Chunk {
	if len(chunks) == 0 || len(terms) == 0 || len(rules) == 0 {
		return chunks
	}
	termSet := map[string]struct{}{}
	for _, t := range terms {
		termSet[t] = struct{}{}
	}
	for i := range chunks {
		for _, rule := range rules {
			if !termsMatch(rule.Terms, termSet) {
				continue
			}
			if pathMatch(rule.Paths, chunks[i].FilePath) {
				chunks[i].Score += rule.Bonus
			}
		}
	}
	return chunks
}

func termsMatch(ruleTerms []string, termSet map[string]struct{}) bool {
	for _, t := range ruleTerms {
		if _, ok := termSet[t]; ok {
			return true
		}
	}
	return false
}

func pathMatch(prefixes []string, path string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
