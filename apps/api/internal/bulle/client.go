package bulle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const anthropicVersion = "2023-06-01"

var (
	// ErrNoKey means Bulle holds no model key for this user or instance.
	ErrNoKey = errors.New("bulle: no model key is set up in Bulle")
	// ErrUnreachable means Bulle did not answer, or answered with an error.
	ErrUnreachable = errors.New("bulle: Bulle is unreachable")
	// ErrUnsupported means Bulle runs a backend this client cannot speak.
	ErrUnsupported = errors.New("bulle: Bulle's model backend is not supported")
)

// Offer is the model setup Bulle runs for one user, persona included.
type Offer struct {
	Backend string   `json:"backend"`
	Model   string   `json:"model"`
	Models  []string `json:"models"`
	Persona string   `json:"persona"`
}

// Prompt is one chat completion request.
type Prompt struct {
	System    string
	User      string
	MaxTokens int
}

// Client calls Bulle as a Courrier user.
type Client struct {
	base   string
	tokens *tokens
	http   *http.Client
}

// New builds a client, or returns nil when BULLE_URL or the Registre client is
// not configured, which turns the assistant off.
func New(cfg Config, ids Identities, httpClient *http.Client) *Client {
	if cfg.BulleURL == "" || cfg.Issuer == "" || cfg.ClientID == "" || ids == nil {
		return nil
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 90 * time.Second}
	}
	return &Client{
		base:   strings.TrimRight(cfg.BulleURL, "/"),
		tokens: &tokens{cfg: cfg, ids: ids, http: httpClient, now: time.Now, cache: map[int64]cached{}},
		http:   httpClient,
	}
}

// Offer reads the backend, model and persona Bulle uses for the user.
func (c *Client) Offer(ctx context.Context, userID int64) (Offer, error) {
	var offer Offer
	err := c.call(ctx, userID, request{method: http.MethodGet, path: "/api/llm"}, &offer)
	return offer, err
}

// Complete runs one completion through Bulle's key-holding proxy and returns
// the model's text.
func (c *Client) Complete(ctx context.Context, userID int64, offer Offer, prompt Prompt) (string, error) {
	model := offer.Model
	if model == "" && len(offer.Models) > 0 {
		model = offer.Models[0]
	}
	if model == "" {
		return "", ErrUnreachable
	}
	switch offer.Backend {
	case "anthropic":
		return c.anthropic(ctx, userID, model, prompt)
	case "openai", "openrouter":
		return c.openAI(ctx, userID, offer.Backend, model, prompt)
	}
	return "", ErrUnsupported
}

func (c *Client) anthropic(ctx context.Context, userID int64, model string, p Prompt) (string, error) {
	body := map[string]any{"model": model, "max_tokens": p.MaxTokens, "system": p.System,
		"messages": []map[string]string{{"role": "user", "content": p.User}}}
	var out struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := c.call(ctx, userID, request{method: http.MethodPost, path: "/api/llm/anthropic/v1/messages", payload: body, headers: map[string]string{"anthropic-version": anthropicVersion}}, &out); err != nil {
		return "", err
	}
	var text strings.Builder
	for _, block := range out.Content {
		if block.Type == "text" {
			text.WriteString(block.Text)
		}
	}
	return nonEmpty(text.String())
}

func (c *Client) openAI(ctx context.Context, userID int64, backend, model string, p Prompt) (string, error) {
	body := map[string]any{"model": model, "max_tokens": p.MaxTokens,
		"messages": []map[string]string{{"role": "system", "content": p.System}, {"role": "user", "content": p.User}}}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := c.call(ctx, userID, request{method: http.MethodPost, path: "/api/llm/" + backend + "/chat/completions", payload: body}, &out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", ErrUnreachable
	}
	return nonEmpty(out.Choices[0].Message.Content)
}

func nonEmpty(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", ErrUnreachable
	}
	return text, nil
}

type request struct {
	method  string
	path    string
	payload any
	headers map[string]string
}

func (c *Client) call(ctx context.Context, userID int64, r request, out any) error {
	res, err := c.send(ctx, userID, r)
	if err == nil && res.StatusCode == http.StatusUnauthorized {
		res.Body.Close()
		c.tokens.forget(userID)
		res, err = c.send(ctx, userID, r)
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(io.LimitReader(res.Body, maxBody))
	switch {
	case res.StatusCode == http.StatusConflict:
		return ErrNoKey
	case res.StatusCode/100 != 2:
		return fmt.Errorf("%w: %s answered %d", ErrUnreachable, r.path, res.StatusCode)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("%w: %v", ErrUnreachable, err)
	}
	return nil
}

func (c *Client) send(ctx context.Context, userID int64, r request) (*http.Response, error) {
	token, err := c.tokens.token(ctx, userID)
	if err != nil {
		return nil, err
	}
	var body io.Reader
	if r.payload != nil {
		data, err := json.Marshal(r.payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, r.method, c.base+r.path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if r.payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnreachable, err)
	}
	return res, nil
}
