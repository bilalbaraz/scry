package ask

import "sort"

// ApplyMatchPreference prefers chunks with >=2 distinct term matches, fallback to >=1.
func ApplyMatchPreference(chunks []Chunk, terms []string) []Chunk {
	if len(chunks) == 0 || len(terms) == 0 {
		return nil
	}
	var strong []Chunk
	var weak []Chunk
	for _, ch := range chunks {
		matches := DistinctTermMatchCount(terms, ch.Text)
		ch.MatchCount = matches
		if matches >= 2 {
			strong = append(strong, ch)
		} else if matches >= 1 {
			weak = append(weak, ch)
		}
	}
	if len(strong) > 0 {
		return strong
	}
	return weak
}

// SortCandidates sorts by score desc, then path asc, then start line asc.
func SortCandidates(chunks []Chunk) []Chunk {
	sort.Slice(chunks, func(i, j int) bool {
		if chunks[i].Score == chunks[j].Score {
			if chunks[i].FilePath == chunks[j].FilePath {
				return chunks[i].StartLine < chunks[j].StartLine
			}
			return chunks[i].FilePath < chunks[j].FilePath
		}
		return chunks[i].Score > chunks[j].Score
	})
	return chunks
}
