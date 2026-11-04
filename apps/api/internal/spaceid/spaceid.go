// Package spaceid canonicalises the space ids clients send. spaces.id is a
// Postgres uuid, which accepts many spellings of one id, while the space_id
// columns of accounts, folders, emails, contacts and templates are text and
// compare byte for byte. Every client-supplied id goes through Optional before
// it is stored or filtered on, so one space never has two spellings.
package spaceid

import (
	"regexp"
	"strings"

	"github.com/FacileStudio/tronc/errors"
)

// canonical is the only spelling stored: lowercase, hyphenated, no braces.
var canonical = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// Optional reads an optional space id: nil or "" is the personal scope and
// answers nil, an uppercase id is lowercased, and any other spelling, such as
// braces or missing hyphens, is refused with 400.
func Optional(raw *string) (*string, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	id := strings.ToLower(*raw)
	if !canonical.MatchString(id) {
		return nil, errors.Invalid("space_id must be a space identifier")
	}
	return &id, nil
}

// Query reads the optional space_id query parameter the same way.
func Query(raw string) (*string, error) {
	return Optional(&raw)
}
