package openai

import (
	"argus-agent/internal/llm"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func mapRequest(req llm.Request) responses.ResponseNewParams {
	return responses.ResponseNewParams{
		Input:           responses.ResponseNewParamsInputUnion{OfString: openai.String(req.Message)},
		Model:           req.Model,
		MaxOutputTokens: openai.Int(req.MaxOutputTokens),
		Temperature:     openai.Float(req.Temperature),
	}
}
