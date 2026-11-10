package openai

import (
	"argus-agent/internal/config"
	"fmt"
	"net/url"
	"strings"

	sdk "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type APIMode string

const (
	APIResponses   APIMode = "responses"
	APICompletions APIMode = "chat_completions"

	OPENAI_HOSTNAME string = "api.openai.com"
)

type Client struct {
	api  *sdk.Client
	mode APIMode
}

func NewClient(cfg config.LLM) *Client {
	client := sdk.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
		option.WithMaxRetries(2),
	)

	return &Client{
		api:  &client,
		mode: parseAPIMode(cfg),
	}
}

func parseAPIMode(cfg config.LLM) APIMode {
	if cfg.APIMode == "" {
		fmt.Println("LLM base url not provided. use default ", APIResponses)
		return APIResponses
	}
	switch APIMode(cfg.APIMode) {
	case APIResponses, APICompletions:
		return APIMode(cfg.APIMode)
	default:
		if cfg.BaseURL == "" {
			cfg.BaseURL = OPENAI_HOSTNAME
			return APIResponses
		}
		u, err := url.Parse(cfg.BaseURL)
		if strings.EqualFold(u.Hostname(), OPENAI_HOSTNAME) || err != nil {
			return APIResponses
		}
		return APICompletions
	}
}
