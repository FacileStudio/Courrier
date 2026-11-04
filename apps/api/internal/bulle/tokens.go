// Package bulle lets Courrier's backend call Bulle as the signed-in user: it
// trades the user's Registre token for a five-minute aud=bulle token (RFC 8693)
// and talks to Bulle's model proxy with it. Courrier never holds a model key.
package bulle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/FacileStudio/porte"
)

const (
	exchangeGrant   = "urn:ietf:params:oauth:grant-type:token-exchange"
	accessTokenType = "urn:ietf:params:oauth:token-type:access_token"
	audience        = "bulle"
	scope           = "bulle.read bulle.write"
	margin          = 30 * time.Second
	maxBody         = 1 << 20
)

// ErrNotLinked means the user has no usable Registre session to trade, or
// Registre refused to trade it for a Bulle token.
var ErrNotLinked = errors.New("bulle: this account is not linked to Bulle")

// Identities is the part of porte's identity store the exchange needs.
type Identities interface {
	ListByUser(ctx context.Context, userID int64) ([]porte.StoredIdentity, error)
	Save(ctx context.Context, identity porte.StoredIdentity) error
}

// Config is where Bulle and Registre live and how Courrier signs in to Registre.
type Config struct {
	BulleURL     string
	Issuer       string
	ClientID     string
	ClientSecret string
}

type cached struct {
	token  string
	expiry time.Time
}

type tokens struct {
	cfg  Config
	ids  Identities
	http *http.Client
	now  func() time.Time

	mu       sync.Mutex
	endpoint string
	cache    map[int64]cached
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Error        string `json:"error"`
}

func (t *tokens) forget(userID int64) {
	t.mu.Lock()
	delete(t.cache, userID)
	t.mu.Unlock()
}

func (t *tokens) token(ctx context.Context, userID int64) (string, error) {
	t.mu.Lock()
	hit, ok := t.cache[userID]
	t.mu.Unlock()
	if ok && t.now().Add(margin).Before(hit.expiry) {
		return hit.token, nil
	}
	subject, err := t.subject(ctx, userID)
	if err != nil {
		return "", err
	}
	out, err := t.post(ctx, url.Values{
		"grant_type": {exchangeGrant}, "subject_token": {subject}, "subject_token_type": {accessTokenType},
		"audience": {audience}, "scope": {scope},
	})
	if err != nil {
		return "", err
	}
	next := cached{token: out.AccessToken, expiry: t.now().Add(time.Duration(out.ExpiresIn) * time.Second)}
	t.mu.Lock()
	t.cache[userID] = next
	t.mu.Unlock()
	return next.token, nil
}

func (t *tokens) subject(ctx context.Context, userID int64) (string, error) {
	ids, err := t.ids.ListByUser(ctx, userID)
	if err != nil {
		return "", err
	}
	for _, id := range ids {
		if id.Provider != t.cfg.Issuer || (id.Tokens.AccessToken == "" && id.Tokens.RefreshToken == "") {
			continue
		}
		if id.Tokens.AccessToken != "" && t.now().Add(margin).Before(id.Tokens.Expiry) {
			return id.Tokens.AccessToken, nil
		}
		return t.refresh(ctx, id)
	}
	return "", ErrNotLinked
}

func (t *tokens) refresh(ctx context.Context, id porte.StoredIdentity) (string, error) {
	if id.Tokens.RefreshToken == "" {
		return "", ErrNotLinked
	}
	out, err := t.post(ctx, url.Values{"grant_type": {"refresh_token"}, "refresh_token": {id.Tokens.RefreshToken}})
	if err != nil {
		return "", err
	}
	id.Tokens.AccessToken = out.AccessToken
	id.Tokens.Expiry = t.now().Add(time.Duration(out.ExpiresIn) * time.Second)
	if out.RefreshToken != "" {
		id.Tokens.RefreshToken = out.RefreshToken
	}
	if err := t.ids.Save(ctx, id); err != nil {
		return "", err
	}
	return out.AccessToken, nil
}

func (t *tokens) post(ctx context.Context, form url.Values) (tokenResponse, error) {
	endpoint, err := t.tokenEndpoint(ctx)
	if err != nil {
		return tokenResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return tokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(url.QueryEscape(t.cfg.ClientID), url.QueryEscape(t.cfg.ClientSecret))
	res, err := t.http.Do(req)
	if err != nil {
		return tokenResponse{}, fmt.Errorf("bulle: registre: %w", err)
	}
	defer res.Body.Close()
	var out tokenResponse
	_ = json.NewDecoder(io.LimitReader(res.Body, maxBody)).Decode(&out)
	switch {
	case res.StatusCode >= 500:
		return out, fmt.Errorf("bulle: registre answered %d", res.StatusCode)
	case res.StatusCode != http.StatusOK || out.AccessToken == "":
		return out, fmt.Errorf("%w: registre said %q", ErrNotLinked, out.Error)
	}
	return out, nil
}

func (t *tokens) tokenEndpoint(ctx context.Context) (string, error) {
	t.mu.Lock()
	known := t.endpoint
	t.mu.Unlock()
	if known != "" {
		return known, nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(t.cfg.Issuer, "/")+"/.well-known/openid-configuration", nil)
	if err != nil {
		return "", err
	}
	res, err := t.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("bulle: registre discovery: %w", err)
	}
	defer res.Body.Close()
	var doc struct {
		TokenEndpoint string `json:"token_endpoint"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, maxBody)).Decode(&doc); err != nil || doc.TokenEndpoint == "" {
		return "", errors.New("bulle: registre discovery returned no token endpoint")
	}
	t.mu.Lock()
	t.endpoint = doc.TokenEndpoint
	t.mu.Unlock()
	return doc.TokenEndpoint, nil
}
