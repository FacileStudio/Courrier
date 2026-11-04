package schemas

import (
	"maps"
	"testing"

	"github.com/FacileStudio/tronc/testdb"
)

// TestNonCanonicalSpaceIDs counts "" and odd spellings, never canonical ids
// or personal rows, and changes nothing.
func TestNonCanonicalSpaceIDs(t *testing.T) {
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	testORM.Exec("DELETE FROM accounts")
	testORM.Exec("DELETE FROM email_templates")
	for _, id := range []string{"", "9EFDB13C-CF46-40A7-9E9F-5781FD4DDBCE", "9efdb13c-cf46-40a7-9e9f-5781fd4ddbce"} {
		spaceID := id
		testORM.Create(&Account{Name: "a", SpaceID: &spaceID})
	}
	testORM.Create(&Account{Name: "personal"})
	brace := "{9efdb13c-cf46-40a7-9e9f-5781fd4ddbce}"
	testORM.Create(&EmailTemplate{Name: "t", SpaceID: &brace})
	found, err := NonCanonicalSpaceIDs(testORM)
	if err != nil {
		t.Fatal(err)
	}
	if want := map[string]int64{"accounts": 2, "email_templates": 1}; !maps.Equal(found, want) {
		t.Fatalf("found %v, want %v", found, want)
	}
	var left int64
	testORM.Model(&Account{}).Where("space_id = ''").Count(&left)
	if left != 1 {
		t.Fatalf("the audit changed rows: %d empty ids left", left)
	}
}
