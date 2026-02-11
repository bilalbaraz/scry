package chunk

type Chunk struct {
	ID        string
	FilePath  string
	StartLine int
	EndLine   int
	Text      string
	Lang      string
}
