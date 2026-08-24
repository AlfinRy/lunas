package config

import "os"

// Config is read from the environment; every value has a demo-safe default so
// the product is fully usable with zero setup (judges: `go run ./cmd/lunas`).
type Config struct {
	Port       string
	DBPath     string
	LLMAPIKey  string
	LLMBaseURL string
	LLMModel   string
}

// TemplateMode is true when no LLM provider is configured: the agent falls
// back to the deterministic chase templates from docs/ux-writing.md.
func (c *Config) TemplateMode() bool { return c.LLMAPIKey == "" }

func Load() *Config {
	return &Config{
		Port:       getenv("LUNAS_PORT", "8080"),
		DBPath:     getenv("LUNAS_DB", "lunas.db"),
		LLMAPIKey:  os.Getenv("LLM_API_KEY"),
		LLMBaseURL: getenv("LLM_BASE_URL", "https://api.openai.com/v1"),
		LLMModel:   getenv("LLM_MODEL", "gpt-4o-mini"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
