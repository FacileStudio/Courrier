package mcp

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// decode reads a tool's structured result into out.
func decode(t *testing.T, result *sdk.CallToolResult, out any) {
	t.Helper()
	if result.IsError {
		t.Fatalf("tool failed: %s", text(result))
	}
	raw, _ := json.Marshal(result.StructuredContent)
	if err := json.Unmarshal(raw, out); err != nil {
		t.Fatal(err)
	}
}

// text is the first text content of a result.
func text(result *sdk.CallToolResult) string {
	for _, content := range result.Content {
		if c, ok := content.(*sdk.TextContent); ok {
			return c.Text
		}
	}
	return ""
}

// TestToolsCarryAnnotationsInNameOrder checks every tool is prefixed,
// annotated with explicit hints, and listed in a stable order.
func TestToolsCarryAnnotationsInNameOrder(t *testing.T) {
	f := newFixture(t)
	listed, err := f.connect(t, f.token(t), clientOptions{}).ListTools(t.Context(), nil)
	if err != nil {
		t.Fatal(err)
	}
	names := []string{}
	for _, tool := range listed.Tools {
		a := tool.Annotations
		if !strings.HasPrefix(tool.Name, "courrier_") || a == nil || a.DestructiveHint == nil || a.OpenWorldHint == nil || a.Title == "" {
			t.Fatalf("tool %s lacks its prefix or explicit annotations: %+v", tool.Name, a)
		}
		names = append(names, tool.Name)
	}
	if len(names) != 7 || !slices.IsSorted(names) {
		t.Fatalf("want the 7 tools sorted by name, got %v", names)
	}
}

// TestSearchPagesWithoutRepeats walks 45 messages in pages of 20.
func TestSearchPagesWithoutRepeats(t *testing.T) {
	f := newFixture(t)
	start := time.Date(2026, 11, 1, 9, 0, 0, 0, time.UTC)
	for i := range 45 {
		f.mail.add(schemas.Email{Subject: fmt.Sprintf("invoice %02d", i), FromAddress: "a@example.test", Date: start.Add(time.Duration(i%9) * time.Hour)})
	}
	session := f.connect(t, f.token(t), clientOptions{})
	seen, cursor, pages := map[int64]bool{}, "", 0
	for {
		var page searchOutput
		decode(t, call(t, session, toolSearch, map[string]any{"query": "invoice", "cursor": cursor}), &page)
		pages++
		for _, message := range page.Messages {
			if seen[message.MessageID] {
				t.Fatalf("message %d came back twice", message.MessageID)
			}
			seen[message.MessageID] = true
		}
		if cursor = page.NextCursor; cursor == "" {
			break
		}
	}
	if len(seen) != 45 || pages != 3 {
		t.Fatalf("want 45 messages over 3 pages, got %d over %d", len(seen), pages)
	}
}

// TestLongThreadIsCutToTheBudget checks a detailed read of big bodies stays
// under the cap and continues on the next page.
func TestLongThreadIsCutToTheBudget(t *testing.T) {
	f := newFixture(t)
	var first schemas.Email
	for i := range 10 {
		email := f.mail.add(schemas.Email{ThreadID: "t1", Subject: "long", Date: time.Now().Add(time.Duration(i) * time.Minute),
			BodyText: strings.Repeat("word ", 5000)})
		if i == 0 {
			first = email
		}
	}
	session := f.connect(t, f.token(t), clientOptions{})
	result := call(t, session, toolReadThread, map[string]any{"account_id": f.accountID, "message_id": first.ID, "response_format": "detailed"})
	var page threadOutput
	decode(t, result, &page)
	if size := encodedSize(page); size > resultBudget || page.NextCursor == "" || len(page.Messages) >= 10 {
		t.Fatalf("want a page under %d bytes with a cursor, got %d bytes, %d messages", resultBudget, size, len(page.Messages))
	}
	var next threadOutput
	decode(t, call(t, session, toolReadThread, map[string]any{"account_id": f.accountID, "message_id": first.ID,
		"response_format": "detailed", "cursor": page.NextCursor}), &next)
	if next.Messages[0].MessageID == page.Messages[len(page.Messages)-1].MessageID {
		t.Fatal("the next page must start after the last message shown")
	}
}

// TestUnknownAccountIsAnActionableError checks a business error is an
// isError result that names the tool to call instead.
func TestUnknownAccountIsAnActionableError(t *testing.T) {
	f := newFixture(t)
	result := call(t, f.connect(t, f.token(t), clientOptions{}), toolSearch, map[string]any{"account_id": 999})
	if !result.IsError || !strings.Contains(text(result), toolListAccounts) {
		t.Fatalf("want an isError result pointing at %s, got %q", toolListAccounts, text(result))
	}
}
