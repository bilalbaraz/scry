package search

import (
	"sort"

	"scry/pkg/index/lexical"
	"scry/pkg/metadata"
)

type Result struct {
	Chunk metadata.ChunkView
	Score float64
}

type Store interface {
	TermHits(term string) ([]metadata.TermHit, error)
	GetChunksByIDs(ids []string) ([]metadata.ChunkView, error)
}

type Engine struct {
	Store Store
}

func New(store Store) *Engine {
	return &Engine{Store: store}
}

func (e *Engine) Search(query string, limit int) ([]Result, error) {
	terms := lexical.Tokenize(query)
	if len(terms) == 0 {
		return nil, nil
	}

	scores := map[string]float64{}
	for _, term := range terms {
		hits, err := e.Store.TermHits(term)
		if err != nil {
			return nil, err
		}
		for _, h := range hits {
			scores[h.ChunkID] += float64(h.TF)
		}
	}
	if len(scores) == 0 {
		return nil, nil
	}

	ids := make([]string, 0, len(scores))
	for id := range scores {
		ids = append(ids, id)
	}
	chunks, err := e.Store.GetChunksByIDs(ids)
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(chunks))
	for _, ch := range chunks {
		results = append(results, Result{Chunk: ch, Score: scores[ch.ID]})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Chunk.FilePath < results[j].Chunk.FilePath
		}
		return results[i].Score > results[j].Score
	})

	if limit > 0 && len(results) > limit {
		return results[:limit], nil
	}
	return results, nil
}
