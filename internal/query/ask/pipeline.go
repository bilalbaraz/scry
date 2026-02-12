package ask

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
	preferred := ApplyMatchPreference(filtered, terms)
	boosted := ApplyBoosts(preferred, terms, DefaultBoostRules())
	boosted = ApplyWhitelistPromotion(boosted, terms, DefaultWhitelistRules())
	boosted = SortCandidates(boosted)
	selected := SelectTopEvidence(boosted, opts.MaxEvidence)
	var evidence []Evidence
	for _, ch := range selected {
		snippet := SnippetAroundTerm(ch.Text, terms, opts.SnippetChars)
		evidence = append(evidence, Evidence{Chunk: ch, Snippet: snippet})
	}
	return evidence
}
