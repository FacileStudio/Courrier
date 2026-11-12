package spaces

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/porte/directory"
	portepg "github.com/FacileStudio/porte/pg"
)

type roundTrip func(*http.Request) (*http.Response, error)

func (f roundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// registreErased is a Registre whose deletion list names subject.
func registreErased(subject string) http.RoundTripper {
	return roundTrip(func(r *http.Request) (*http.Response, error) {
		recorder := httptest.NewRecorder()
		if r.URL.Path != "/api/directory/deletions" {
			recorder.WriteHeader(http.StatusNotFound)
			return recorder.Result(), nil
		}
		recorder.Header().Set("Content-Type", "application/json")
		body := `{"deleted":[{"sub":"` + subject + `","deleted_at":"` + time.Now().UTC().Format(time.RFC3339Nano) + `"}]}`
		if _, err := recorder.WriteString(body); err != nil {
			return nil, err
		}
		return recorder.Result(), nil
	})
}

// seedCredentials gives a person a session, an API token, a push device and
// an open MCP action, the rows that hang off the person rather than a mailbox.
func seedCredentials(t *testing.T, userID int64) {
	t.Helper()
	for _, statement := range []string{
		`INSERT INTO porte_sessions (token_hash, user_id, created_at, last_used_at, expires_at)
		 VALUES ('session-' || ?, ?::bigint, now(), now(), now() + interval '1 day')`,
		`INSERT INTO porte_sessions (token_hash, user_id, label, created_at, last_used_at)
		 VALUES ('token-' || ?, ?::bigint, 'cli', now(), now())`,
		`INSERT INTO push_devices (kind, token, user_id) VALUES ('fcm', 'device-' || ?, ?::bigint)`,
		`INSERT INTO mcp_actions (id, user_id, kind, target, digest, expires_at)
		 VALUES ('action-' || ?, ?, 'send', '1/1', 'd', now() + interval '1 hour')`,
	} {
		id := strconv.FormatInt(userID, 10)
		if err := testORM.Exec(statement, id, id).Error; err != nil {
			t.Fatalf("seed credentials: %v", err)
		}
	}
}

// referencing counts, per table with a user_id column, the rows that still
// point at userID. Rows that belong to a space are not the person's.
func referencing(t *testing.T, userID int64) map[string]int64 {
	t.Helper()
	var tables []string
	err := testORM.Raw(`SELECT table_name FROM information_schema.columns
		WHERE table_schema = current_schema() AND column_name = 'user_id' ORDER BY table_name`).Scan(&tables).Error
	if err != nil {
		t.Fatalf("tables: %v", err)
	}
	out := map[string]int64{}
	for _, table := range tables {
		query := `SELECT count(*) FROM ` + table + ` WHERE user_id::text = ?`
		if table == "accounts" || table == "email_templates" {
			query += ` AND space_id IS NULL`
		}
		if n := count(t, testORM, query, strconv.FormatInt(userID, 10)); n > 0 {
			out[table] = n
		}
	}
	return out
}

// newErasureDeletions builds the receiver as main.go does: porte/kit's
// EraseAccount, porte's identity eraser and holder check, and a Registre that lists sub-gone.
func newErasureDeletions(t *testing.T) *Deletions {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sqlDB, err := testORM.DB()
	if err != nil {
		t.Fatal(err)
	}
	users, err := auth.NewUserStore(context.Background(), sqlDB, logger)
	if err != nil {
		t.Fatalf("store: %v", err)
	}
	return NewDeletions(testORM, DeletionsConfig{
		Client:        directory.Client{Issuer: testIssuer, HTTP: &http.Client{Transport: registreErased("sub-gone")}},
		EraseAccount:  users.EraseAccount,
		EraseIdentity: portepg.New(sqlDB).EraseIdentity,
		HoldsIdentity: portepg.New(sqlDB).HoldsIdentity,
		StorageDir:    t.TempDir(),
		Logger:        logger,
	})
}

// pollErased runs one deletion poll from a clean receiver state, the path
// the background poller takes.
func pollErased(t *testing.T, deletions *Deletions) {
	t.Helper()
	if err := testORM.Exec(`DELETE FROM directory_deletion_state`).Error; err != nil {
		t.Fatal(err)
	}
	if err := deletions.receiver.Poll(context.Background()); err != nil {
		t.Fatalf("poll: %v", err)
	}
}

// After Registre erases a person, no row in any table with a user_id column
// references them, except the mailboxes and templates that belong to a space.
// Another person keeps every kind of row.
func TestErasureLeavesNoRowReferencingThePerson(t *testing.T) {
	f := seedErasure(t)
	pollErased(t, f.deletions)

	if left := referencing(t, f.gone); len(left) > 0 {
		t.Errorf("rows still reference the erased person: %v", left)
	}
	if left := referencing(t, f.other); len(left) != 6 {
		t.Errorf("another person's rows = %v, want identity, sessions, push, mcp, accounts and templates", left)
	}
}
