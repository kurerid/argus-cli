package openai

import (
	"argus-agent/internal/llm"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

func mapRequestForResponses(req llm.Request) responses.ResponseNewParams {
	return responses.ResponseNewParams{
		Input:           responses.ResponseNewParamsInputUnion{OfString: openai.String(req.Message)},
		Model:           req.Model,
		MaxOutputTokens: openai.Int(req.MaxOutputTokens),
		Temperature:     openai.Float(req.Temperature),
	}
}

func mapRequestForCompletions(req llm.Request) openai.ChatCompletionNewParams {
	return openai.ChatCompletionNewParams{
		Model: req.Model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfUser: &openai.ChatCompletionUserMessageParam{Content: openai.ChatCompletionUserMessageParamContentUnion{OfString: openai.String(req.Message)}},
			},
		},
		MaxTokens:   openai.Int(req.MaxOutputTokens),
		Temperature: openai.Float(req.Temperature),
	}
}
