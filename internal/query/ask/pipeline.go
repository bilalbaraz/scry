package ask

import (
	"sort"
)

type Options struct {
	MaxEvidence  int
	SnippetChars int
}

type Evidence struct {
	Chunk   Chunk
	Snippet string
}

func BuildEvidence(chunks []Chunk, terms []string, opts Options) []Evidence {
	filtered := FilterByQueryTerms(chunks, terms)
	boosted := ApplyBoosts(filtered, terms, DefaultBoostRules())
	sort.Slice(boosted, func(i, j int) bool {
		if boosted[i].Score == boosted[j].Score {
			if boosted[i].FilePath == boosted[j].FilePath {
				return boosted[i].ID < boosted[j].ID
			}
			return boosted[i].FilePath < boosted[j].FilePath
		}
		return boosted[i].Score > boosted[j].Score
	})
	selected := SelectTopEvidence(boosted, opts.MaxEvidence)
	var evidence []Evidence
	for _, ch := range selected {
		snippet := SnippetAroundTerm(ch.Text, terms, opts.SnippetChars)
		evidence = append(evidence, Evidence{Chunk: ch, Snippet: snippet})
	}
	return evidence
}
