package lexical

import (
	"strings"
	"unicode"
)

type Posting struct {
	Term    string
	ChunkID string
	TF      int
}

type InvertedIndex struct {
	Postings map[string][]Posting
}

func New() *InvertedIndex {
	return &InvertedIndex{Postings: make(map[string][]Posting)}
}

func (idx *InvertedIndex) Add(chunkID string, text string) []Posting {
	terms := tokenize(text)
	counts := map[string]int{}
	for _, term := range terms {
		counts[term]++
	}
	var postings []Posting
	for term, tf := range counts {
		p := Posting{Term: term, ChunkID: chunkID, TF: tf}
		idx.Postings[term] = append(idx.Postings[term], p)
		postings = append(postings, p)
	}
	return postings
}

func tokenize(s string) []string {
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
