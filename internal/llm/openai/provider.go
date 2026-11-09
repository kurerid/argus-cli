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

	resp, err := p.client.api.Responses.New(ctx, mapRequest(req))
	if err != nil {
		fmt.Println("FAILED TO GET RESPONSE FROM MODEL")
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
