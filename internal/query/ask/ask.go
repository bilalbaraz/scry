package ask

import "fmt"

type Decision struct {
	Evidence []Evidence
	Reason   string
}

type AskOptions struct {
	MaxEvidence  int
	SnippetChars int
	MinScore     float64
}

func Decide(chunks []Chunk, terms []string, opts AskOptions) Decision {
	evidence := BuildEvidence(chunks, terms, Options{MaxEvidence: opts.MaxEvidence, SnippetChars: opts.SnippetChars})
	if len(evidence) == 0 {
		return Decision{Reason: "no_evidence"}
	}
	if totalScore(evidence) < opts.MinScore {
		return Decision{Reason: "low_score"}
	}
	return Decision{Evidence: evidence}
}

func totalScore(evidence []Evidence) float64 {
	var sum float64
	for _, e := range evidence {
		sum += e.Chunk.Score
	}
	return sum
}

func AnswerHeader(evidence []Evidence) string {
	return fmt.Sprintf("Found %d relevant evidence chunk(s).", len(evidence))
}
