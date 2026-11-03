package spaces

import (
	"context"

	"github.com/FacileStudio/porte/directory"
)

// importQuery lists every membership of every space with the member's
// Registre subject when Courrier has one. Courrier has no personal spaces, so
// every space is a shared one; a space nobody belongs to is left out.
const importQuery = `
SELECT s.id AS space_id, s.name, m.role, u.email, coalesce(i.subject, '') AS subject
  FROM spaces s
  JOIN space_members m ON m.space_id = s.id
  JOIN users u ON u.id = m.user_id
  LEFT JOIN porte_identities i ON i.user_id = u.id AND i.provider = ?
 ORDER BY s.created_at, s.id, m.joined_at`

// importRow is one line of importQuery.
type importRow struct {
	SpaceID string
	Name    string
	Role    string
	Email   string
	Subject string
}

// imports is what Courrier hands Registre the first time: each space by name,
// each member by subject, or by email when they never signed in through
// Registre.
func (d *Directory) imports(ctx context.Context) ([]directory.ImportSpace, error) {
	var rows []importRow
	if err := d.orm.WithContext(ctx).Raw(importQuery, d.issuer).Scan(&rows).Error; err != nil {
		return nil, err
	}
	spaces := []directory.ImportSpace{}
	index := map[string]int{}
	for _, row := range rows {
		at, seen := index[row.SpaceID]
		if !seen {
			at = len(spaces)
			index[row.SpaceID] = at
			spaces = append(spaces, directory.ImportSpace{Name: row.Name})
		}
		spaces[at].Members = append(spaces[at].Members, importMember(row))
	}
	return spaces, nil
}

// importMember names a member by subject when there is one, else by email.
func importMember(row importRow) directory.ImportMember {
	if row.Subject != "" {
		return directory.ImportMember{Subject: row.Subject, Role: localRole(row.Role)}
	}
	return directory.ImportMember{Email: row.Email, Role: localRole(row.Role)}
}
