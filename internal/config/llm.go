package config

type LLM struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
	APIMode string `json:"api_mode"`
}
