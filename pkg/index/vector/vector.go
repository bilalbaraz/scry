package vector

type Provider interface {
	Embed(texts []string) ([][]float32, error)
	Dimension() int
	Name() string
	OfflineOnly() bool
}

type Index struct {
	Provider Provider
}

func New(provider Provider) *Index {
	return &Index{Provider: provider}
}
