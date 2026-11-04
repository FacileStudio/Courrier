package schemas

import (
	"fmt"

	"gorm.io/gorm"
)

// spaceIDTables are the tables whose space_id is text, compared byte for byte
// against the canonical id spaces.id prints.
var spaceIDTables = []string{"accounts", "folders", "emails", "contacts", "email_templates"}

// NonCanonicalSpaceIDs counts, per table, the rows whose space_id is set but
// is not a lowercase hyphenated uuid, such as "" or an uppercase spelling
// stored before ids were canonicalised. It only reads: the owner decides what
// to do with those rows.
func NonCanonicalSpaceIDs(db *gorm.DB) (map[string]int64, error) {
	found := map[string]int64{}
	for _, table := range spaceIDTables {
		var count int64
		query := fmt.Sprintf(`SELECT count(*) FROM %s WHERE space_id IS NOT NULL
			AND space_id !~ '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'`, table)
		if err := db.Raw(query).Scan(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			found[table] = count
		}
	}
	return found, nil
}
