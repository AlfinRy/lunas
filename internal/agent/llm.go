package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// DraftProvider shapes chase words. Two implementations ship:
// TemplateProvider (deterministic, always available) and OpenAIProvider
// (any OpenAI-compatible endpoint). The policy engine is unaffected by which
// one answers — decisions are not delegated to the model.
type DraftProvider interface {
	Name() string
	Draft(ctx context.Context, f ChaseFacts) (DraftEmail, error)
}

// TemplateProvider always answers; it is the product's floor.
type TemplateProvider struct{}

func (TemplateProvider) Name() string { return "template" }
func (TemplateProvider) Draft(_ context.Context, f ChaseFacts) (DraftEmail, error) {
	return Template(f), nil
}

// OpenAIProvider talks to any OpenAI-compatible /chat/completions endpoint.
type OpenAIProvider struct {
	BaseURL string // e.g. https://api.openai.com/v1
	APIKey  string
	Model   string
	Client  *http.Client
}

func (p *OpenAIProvider) Name() string { return "llm:" + p.Model }

func (p *OpenAIProvider) Draft(ctx context.Context, f ChaseFacts) (DraftEmail, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	tmpl := Template(f) // grounding: the canonical copy this stage must respect
	sys := fmt.Sprintf(`You draft invoice follow-up emails for Lunas, a collections assistant, on behalf of %s.
Rules:
- Stage %s (%s tone). Follow the reference draft's structure, facts, and firmness exactly; improve flow only.
- Keep every number, date, invoice number, and the payment link verbatim.
- Never invent fees, deadlines, or claims not present in the reference.
- Sign as %s. Keep the automated-footer line at the end.
- Reply with JSON only: {"subject": string, "body": string}.`,
		f.SenderName, f.Stage, f.Stage.Tone(), f.SenderName)

	user := fmt.Sprintf("Reference draft:\nSubject: %s\n\n%s", tmpl.Subject, tmpl.Body)

	payload := map[string]any{
		"model": p.Model,
		"messages": []map[string]string{
			{"role": "system", "content": sys},
			{"role": "user", "content": user},
		},
		"temperature":       0.4,
		"response_format":   map[string]string{"type": "json_object"},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return DraftEmail{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		strings.TrimSuffix(p.BaseURL, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return DraftEmail{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

	res, err := (p.ClientOrDefault()).Do(req)
	if err != nil {
		return DraftEmail{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return DraftEmail{}, fmt.Errorf("provider status %d", res.StatusCode)
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return DraftEmail{}, err
	}
	if len(out.Choices) == 0 {
		return DraftEmail{}, fmt.Errorf("empty completion")
	}
	var parsed struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	if err := json.Unmarshal([]byte(out.Choices[0].Message.Content), &parsed); err != nil {
		return DraftEmail{}, err
	}
	if parsed.Subject == "" || parsed.Body == "" {
		return DraftEmail{}, fmt.Errorf("incomplete draft")
	}
	return DraftEmail{Subject: parsed.Subject, Body: parsed.Body}, nil
}

func (p *OpenAIProvider) ClientOrDefault() *http.Client {
	if p.Client != nil {
		return p.Client
	}
	return http.DefaultClient
}

// ResilientProvider wraps a primary with the template floor: on any primary
// error the canonical template answers instead. The product never blocks on
// the LLM.
type ResilientProvider struct {
	Primary DraftProvider
	Floor   DraftProvider
}

func (r ResilientProvider) Name() string {
	if r.Primary != nil {
		return r.Primary.Name()
	}
	return "template"
}

func (r ResilientProvider) Draft(ctx context.Context, f ChaseFacts) (DraftEmail, error) {
	if r.Primary != nil {
		if d, err := r.Primary.Draft(ctx, f); err == nil {
			return d, nil
		}
	}
	return r.Floor.Draft(ctx, f)
}
