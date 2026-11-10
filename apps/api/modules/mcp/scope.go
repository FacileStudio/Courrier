package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/FacileStudio/porte/mcpauth"
	"github.com/modelcontextprotocol/go-sdk/auth"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// rpcCall is the part of a JSON-RPC message the write gate reads.
type rpcCall struct {
	Method string `json:"method"`
	Params struct {
		Name string `json:"name"`
	} `json:"params"`
}

// writeGate answers 403 insufficient_scope when a token without the write
// scope calls a tool in writes, so the client can ask the person for more
// access (step-up) instead of failing inside the tool. It reads the body
// rather than the Mcp-Name header, because the body is what the transport
// acts on.
func writeGate(writes map[string]bool, scope, metadataURL string) func(http.Handler) http.Handler {
	challenge := fmt.Sprintf(`Bearer error="insufficient_scope", scope=%q, resource_metadata=%q`, scope, metadataURL)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if mcpauth.HasScope(auth.TokenInfoFromContext(r.Context()), scope) || r.Method != http.MethodPost {
				next.ServeHTTP(w, r)
				return
			}
			body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, sdk.DefaultMaxRequestBodyBytes))
			if err != nil {
				http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))
			if callsAny(body, writes) {
				w.Header().Set("WWW-Authenticate", challenge)
				http.Error(w, "Forbidden: this tool needs the "+scope+" scope", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// callsAny reports whether a JSON-RPC message, or any message of a batch,
// calls one of the named tools.
func callsAny(body []byte, names map[string]bool) bool {
	var calls []rpcCall
	if err := json.Unmarshal(body, &calls); err != nil {
		var call rpcCall
		if json.Unmarshal(body, &call) != nil {
			return false
		}
		calls = []rpcCall{call}
	}
	for _, call := range calls {
		if call.Method == "tools/call" && names[call.Params.Name] {
			return true
		}
	}
	return false
}
