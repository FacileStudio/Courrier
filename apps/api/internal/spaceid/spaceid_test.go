package spaceid

import "testing"

// TestOptional accepts one spelling per id and the personal scope.
func TestOptional(t *testing.T) {
	const id = "9efdb13c-cf46-40a7-9e9f-5781fd4ddbce"
	for _, raw := range []string{id, "9EFDB13C-CF46-40A7-9E9F-5781FD4DDBCE"} {
		got, err := Query(raw)
		if err != nil || got == nil || *got != id {
			t.Errorf("Query(%q) = %v, %v", raw, got, err)
		}
	}
	for _, raw := range []string{"{" + id + "}", "9efdb13ccf4640a79e9f5781fd4ddbce", " " + id, id + "\n", "urn:uuid:" + id, "x"} {
		if _, err := Query(raw); err == nil {
			t.Errorf("Query(%q) was accepted", raw)
		}
	}
	if got, err := Optional(nil); got != nil || err != nil {
		t.Errorf("Optional(nil) = %v, %v", got, err)
	}
	if got, err := Query(""); got != nil || err != nil {
		t.Errorf(`Query("") = %v, %v`, got, err)
	}
}
