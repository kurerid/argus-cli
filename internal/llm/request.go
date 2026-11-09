package llm

type Request struct {
	Model           string
	Message         string
	Temperature     float64
	MaxOutputTokens int64
}
