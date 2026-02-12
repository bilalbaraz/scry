package ask

type Chunk struct {
	ID        string
	FilePath  string
	StartLine int
	EndLine   int
	Text      string
	Score     float64
}
