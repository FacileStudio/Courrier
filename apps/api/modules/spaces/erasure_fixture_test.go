package spaces

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

// newErasureDeletions builds the receiver as main.go does: porte/kit's
// EraseAccount, porte's identity eraser, and a Registre that lists sub-gone.
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
