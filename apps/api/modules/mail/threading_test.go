package mail

import "testing"

func TestParseRefIDs(t *testing.T) {
	got := parseRefIDs("<Root@x> <M2@x>\r\n <Parent@x>")
	want := []string{"root@x", "m2@x", "parent@x"}
	if len(got) != len(want) {
		t.Fatalf("parseRefIDs len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("parseRefIDs[%d] = %q, want %q", i, got[i], want[i])
		}
	}

	// Bracketless fallback for sloppy senders.
	if got := parseRefIDs("a@x, b@x"); len(got) != 2 || got[0] != "a@x" || got[1] != "b@x" {
		t.Errorf("bracketless parse = %v", got)
	}
	if got := parseRefIDs("   "); got != nil {
		t.Errorf("blank parse = %v, want nil", got)
	}
}

func TestCollectThreadIDs(t *testing.T) {
	// Own id, references chain, and in-reply-to deduped into one ordered set.
	got := collectThreadIDs("<C@x>", "<Parent@x>", "<Root@x> <Parent@x>")
	want := []string{"c@x", "root@x", "parent@x"}
	if len(got) != len(want) {
		t.Fatalf("collect len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("collect[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestNewThreadID(t *testing.T) {
	// References root wins.
	if got := newThreadID("<C@x>", "<Parent@x>", "<Root@x> <Parent@x>"); got != "root@x" {
		t.Errorf("newThreadID with refs = %q, want root@x", got)
	}
	// No references: fall back to the parent.
	if got := newThreadID("<C@x>", "<Parent@x>", ""); got != "parent@x" {
		t.Errorf("newThreadID with irt = %q, want parent@x", got)
	}
	// Top-level message threads on its own id.
	if got := newThreadID("<Solo@x>", "", ""); got != "solo@x" {
		t.Errorf("newThreadID solo = %q, want solo@x", got)
	}
}
