package ask

import "strings"

type WhitelistRule struct {
	Terms []string
	Paths []string
	Boost float64
}

func DefaultWhitelistRules() []WhitelistRule {
	return []WhitelistRule{
		{
			Terms: []string{"scan", "ignore", "gitignore", "exclude", "pattern"},
			Paths: []string{"pkg/scan/", "pkg/ignore/"},
			Boost: 1.5,
		},
	}
}

// ApplyWhitelistPromotion adds a soft boost for whitelisted paths when query terms match.
func ApplyWhitelistPromotion(chunks []Chunk, terms []string, rules []WhitelistRule) []Chunk {
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
				chunks[i].Score += rule.Boost
			}
		}
	}
	return chunks
}

func isScanRelated(terms []string) bool {
	for _, t := range terms {
		switch strings.ToLower(t) {
		case "scan", "ignore", "gitignore", "exclude", "pattern":
			return true
		}
	}
	return false
}
