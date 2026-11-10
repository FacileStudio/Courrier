package mcp

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"unicode/utf8"
)

// resultBudget caps the encoded size of one tool result at about 25,000
// tokens, counting three bytes a token for JSON.
const resultBudget = 75_000

// Response formats a read tool accepts.
const (
	formatConcise  = "concise"
	formatDetailed = "detailed"
)

// errCursor is what a model sees when it sends a cursor this server did not
// mint, or one for another tool.
var errCursor = errors.New("the cursor is not valid. Omit it to start again from the first page")

// detailed reports whether the caller asked for the detailed format, and
// rejects any value other than the two formats.
func detailed(format string) (bool, error) {
	switch format {
	case "", formatConcise:
		return false, nil
	case formatDetailed:
		return true, nil
	default:
		return false, errors.New(`response_format must be "concise" or "detailed"`)
	}
}

// encodeCursor makes an opaque cursor out of a position.
func encodeCursor(position any) string {
	raw, _ := json.Marshal(position)
	return base64.RawURLEncoding.EncodeToString(raw)
}

// decodeCursor reads a cursor made by encodeCursor into position. An empty
// cursor leaves position untouched.
func decodeCursor(cursor string, position any) error {
	if cursor == "" {
		return nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || json.Unmarshal(raw, position) != nil {
		return errCursor
	}
	return nil
}

// fitCount returns the largest n in [0, total] for which size(n), the
// encoded size of a result holding the first n items, stays within budget.
// Size grows with n, so a binary search finds it in log(total) encodings.
func fitCount(total int, size func(n int) int) int {
	return max(sort.Search(total+1, func(n int) bool { return size(n) > resultBudget })-1, 0)
}

// encodedSize is the JSON size of v.
func encodedSize(v any) int {
	raw, _ := json.Marshal(v)
	return len(raw)
}

// clip cuts text to at most limit bytes on a rune boundary and reports
// whether anything was cut.
func clip(text string, limit int) (string, bool) {
	if len(text) <= limit {
		return text, false
	}
	cut := limit
	for cut > 0 && !utf8.RuneStart(text[cut]) {
		cut--
	}
	return strings.TrimSpace(text[:cut]), true
}
