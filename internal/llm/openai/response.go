package openai

import "argus-agent/internal/llm"

func mapResponse(responseText string, responseStatus string) *llm.Response {
	return &llm.Response{
		Content:      responseText,
		FinishReason: responseStatus,
	}
}
