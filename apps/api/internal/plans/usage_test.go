package plans

import "testing"

func usageOf(t *testing.T, gate *Gate, userID int64, spaceID *string) Usage {
	t.Helper()
	usage, err := gate.MailboxUsage(testORM, userID, spaceID)
	if err != nil {
		t.Fatal(err)
	}
	return usage
}

// TestMailboxUsageReadsEachScope reports the accounts and ceiling of the
// personal scope and of a space, unlimited before the first snapshot, and a
// ceiling of 0 for a space Registre does not know.
func TestMailboxUsageReadsEachScope(t *testing.T) {
	requireDB(t)
	gate := newGate()
	alice := newUser(t, "alice@facile.studio", "alice")
	pro := newSpace(t, "Pro", "spc_pro")
	local := newSpace(t, "Local", "")
	addAccounts(t, alice, &pro, 1)
	addAccounts(t, alice, nil, 2)

	if got := usageOf(t, gate, alice, nil); got.Used != 2 || got.Limit != -1 {
		t.Fatalf("before the first snapshot: %+v", got)
	}
	var nilGate *Gate
	if got := usageOf(t, nilGate, alice, &pro); got.Used != 1 || got.Limit != -1 {
		t.Fatalf("nil gate: %+v", got)
	}

	apply(t, gate, entitled(1, map[string]int{"spc_pro": 4}, map[string]int{"alice": 3}))
	if got := usageOf(t, gate, alice, nil); got.Used != 2 || got.Limit != 3 {
		t.Fatalf("personal: %+v", got)
	}
	if got := usageOf(t, gate, alice, &pro); got.Used != 1 || got.Limit != 4 {
		t.Fatalf("space: %+v", got)
	}
	if got := usageOf(t, gate, alice, &local); got.Used != 0 || got.Limit != 0 {
		t.Fatalf("unlinked space: %+v", got)
	}
}
