package openai

import (
	"argus-agent/internal/llm"
	"context"
	"errors"
	"fmt"

	"github.com/openai/openai-go/v3/responses"
)

type Provider struct {
	client *Client
}

func NewProvider(client *Client) *Provider {
	return &Provider{
		client: client,
	}
}

func (p *Provider) Generate(
	ctx context.Context,
	req llm.Request,
) (*llm.Response, error) {
	switch p.client.mode {
	case APIResponses:
		return p.useResponsesAPI(ctx, req)
	case APICompletions:
		return p.useCompletionsAPI(ctx, req)
	default:
		return nil, fmt.Errorf("unexpected LLM client mode")
	}
}

func (p *Provider) useResponsesAPI(ctx context.Context, req llm.Request) (*llm.Response, error) {
	resp, err := p.client.api.Responses.New(ctx, mapRequestForResponses(req))
	if err != nil {
		return nil, errors.New("failed to get response from model: " + err.Error())
	}
	switch resp.Status {
	case responses.ResponseStatusCompleted:
		return mapResponse(resp.OutputText(), string(resp.Status)), nil
	case responses.ResponseStatusFailed:
		return nil, fmt.Errorf(
			"model generation failed: %v",
			resp.Error,
		)
	case responses.ResponseStatusIncomplete:
		return nil, fmt.Errorf(
			"model generation incomplete: %v",
			resp.IncompleteDetails,
		)
	default:
		return nil, fmt.Errorf(
			"unexpected response status: %s",
			resp.Status,
		)
	}
}
func (p *Provider) useCompletionsAPI(ctx context.Context, req llm.Request) (*llm.Response, error) {
	resp, err := p.client.api.Chat.Completions.New(ctx, mapRequestForCompletions(req))
	if err != nil {
		fmt.Println("FAILED TO GET RESPONSE FROM MODEL")
		return nil, fmt.Errorf("failed to get response from model: %w", err)
	}
	if len(resp.Choices) > 0 {
		outputText := resp.Choices[0].Message.Content
		return mapResponse(outputText, string(responses.ResponseStatusCompleted)), nil
	}
	return nil, fmt.Errorf("LLM send response with choices len 0")
}
