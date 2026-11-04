package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	jevURL     = "https://api.typesafe.ai/v1/systemone"
	jevModel   = "jev-latest"
	jevTimeout = 20 * time.Second
)

type jevQuestion struct {
	Type         string            `json:"type"`
	Instructions string            `json:"instructions"`
	Criteria     map[string]string `json:"criteria"`
	Options      []string          `json:"options"`
}

type jevAnswer struct {
	Choice     string  `json:"choice"`
	Confidence float64 `json:"confidence"`
}

// jevClient asks TypeSafe's jev decision model one choice question. It is the
// System One wire format Bulle's core/jev speaks, kept to the one call
// Courrier makes so the API module does not pull Bulle's module graph.
type jevClient struct {
	url    string
	apiKey string
	http   *http.Client
}

func jevFromEnv() *jevClient {
	key := strings.TrimSpace(os.Getenv("TYPESAFE_API_KEY"))
	if key == "" {
		return nil
	}
	url := strings.TrimSpace(os.Getenv("TYPESAFE_URL"))
	if url == "" {
		url = jevURL
	}
	return &jevClient{url: url, apiKey: key, http: &http.Client{Timeout: jevTimeout}}
}

func (c *jevClient) choose(ctx context.Context, state any, name string, q jevQuestion) (jevAnswer, error) {
	payload, err := json.Marshal(map[string]any{"state": state, "model": jevModel, "questions": map[string]jevQuestion{name: q}})
	if err != nil {
		return jevAnswer{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(payload))
	if err != nil {
		return jevAnswer{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	res, err := c.http.Do(req)
	if err != nil {
		return jevAnswer{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return jevAnswer{}, err
	}
	if res.StatusCode/100 != 2 {
		return jevAnswer{}, fmt.Errorf("jev: %d %s", res.StatusCode, strings.TrimSpace(string(body)))
	}
	var out struct {
		Answers map[string]jevAnswer `json:"answers"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return jevAnswer{}, err
	}
	return out.Answers[name], nil
}
