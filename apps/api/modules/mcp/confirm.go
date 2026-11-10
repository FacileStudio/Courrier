package mcp

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// confirmInput is the key of the one input request a sensitive tool makes.
const confirmInput = "confirm"

// mrtrVersion is the first protocol version with multi round-trip requests.
// Older clients reach this stateless server without a session, so the SDK's
// fallback (a server-to-client elicitation) cannot run for them.
const mrtrVersion = "2026-07-28"

// verdict is what the person answered to a confirmation.
type verdict int

const (
	unanswered verdict = iota
	accepted
	declined
)

// errConfirmation is the refusal for a confirmation state this server did
// not sign, signed for someone else, or signed too long ago.
var errConfirmation = errors.New("this confirmation is no longer valid. Call the tool again to ask the user afresh")

// confirmSchema is the form the client shows: a single checkbox, ticked.
var confirmSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"confirm": map[string]any{"type": "boolean", "title": "Confirm", "default": true},
	},
	"required": []string{"confirm"},
}

// confirmer asks a client to confirm a sensitive action through MRTR
// input_required, and reads the answer on the retry. The request state it
// hands out names a pending action, the person and an expiry, signed with
// HMAC-SHA256: a client cannot forge a confirmation for an action the server
// did not ask about, and the pending action makes each one single use.
type confirmer struct {
	key []byte
	ttl time.Duration
}

// canElicit reports whether the client asked for form elicitation on a
// protocol version that lets a stateless server use it.
func canElicit(req *sdk.CallToolRequest) bool {
	if req.ProtocolVersion() < mrtrVersion {
		return false
	}
	caps := req.ClientCapabilities()
	if caps == nil || caps.Elicitation == nil {
		return false
	}
	return caps.Elicitation.Form != nil || caps.Elicitation.URL == nil
}

// retrying reports whether the call is the retry of an input_required
// result.
func retrying(req *sdk.CallToolRequest) bool {
	return req.Params.RequestState != "" || len(req.Params.InputResponses) > 0
}

// ask returns the input_required result for pending action id.
func (c confirmer) ask(userID int64, id, message string) *sdk.CallToolResult {
	return &sdk.CallToolResult{
		InputRequests: sdk.InputRequestMap{confirmInput: &sdk.ElicitParams{
			Mode: "form", Message: message, RequestedSchema: confirmSchema,
		}},
		RequestState: c.sign(userID, id, time.Now().Add(c.ttl)),
	}
}

// answer reads the retry of an ask: the pending action it was about and
// whether the person accepted.
func (c confirmer) answer(req *sdk.CallToolRequest, userID int64) (string, verdict, error) {
	id, err := c.verify(req.Params.RequestState, userID, time.Now())
	if err != nil {
		return "", unanswered, err
	}
	result, ok := req.Params.InputResponses[confirmInput].(*sdk.ElicitResult)
	if !ok {
		return "", unanswered, errConfirmation
	}
	if result.Action != "accept" || result.Content["confirm"] == false {
		return id, declined, nil
	}
	return id, accepted, nil
}

// sign makes the request state for pending action id.
func (c confirmer) sign(userID int64, id string, expires time.Time) string {
	payload := strconv.FormatInt(userID, 10) + "." + id + "." + strconv.FormatInt(expires.Unix(), 10)
	return payload + "." + c.mac(payload)
}

// verify returns the pending action id a state names, if the state is this
// server's, for userID, and not expired at now.
func (c confirmer) verify(state string, userID int64, now time.Time) (string, error) {
	parts := strings.Split(state, ".")
	if len(parts) != 4 || parts[0] != strconv.FormatInt(userID, 10) {
		return "", errConfirmation
	}
	payload := strings.Join(parts[:3], ".")
	if !hmac.Equal([]byte(parts[3]), []byte(c.mac(payload))) {
		return "", errConfirmation
	}
	expires, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || now.Unix() > expires {
		return "", errConfirmation
	}
	return parts[1], nil
}

// mac is the hex HMAC-SHA256 of payload.
func (c confirmer) mac(payload string) string {
	sum := hmac.New(sha256.New, c.key)
	sum.Write([]byte(payload))
	return hex.EncodeToString(sum.Sum(nil))
}
