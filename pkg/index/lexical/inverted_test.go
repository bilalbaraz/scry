package lexical

import "testing"

func TestInvertedIndexAdd(t *testing.T) {
	idx := New()
	postings := idx.Add("c1", "alpha beta beta")
	if len(postings) != 2 {
		t.Fatalf("expected 2 postings, got %d", len(postings))
	}

	tf := map[string]int{}
	for _, p := range postings {
		tf[p.Term] = p.TF
		if p.ChunkID != "c1" {
			t.Fatalf("unexpected chunk id: %s", p.ChunkID)
		}
	}
	if tf["alpha"] != 1 || tf["beta"] != 2 {
		t.Fatalf("unexpected tf map: %v", tf)
	}

	if len(idx.Postings["alpha"]) != 1 || len(idx.Postings["beta"]) != 1 {
		t.Fatalf("expected postings to be stored in index")
	}
}
