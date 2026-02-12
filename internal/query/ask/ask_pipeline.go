package ask

// BuildPipeline converts raw chunks into scored candidates, then applies filtering/boosting.
func BuildPipeline(chunks []Chunk, terms []string, opts AskOptions) Decision {
	return Decide(chunks, terms, opts)
}
