package spaces

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/FacileStudio/porte/directory"
)

// fakeRegistre answers 404 to the first refusals imports, as a Registre that
// is not deployed yet does, then accepts. It records the order of calls.
type fakeRegistre struct {
	mu       sync.Mutex
	refusals int
	calls    []string
}

func (f *fakeRegistre) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, r.URL.Path)
	if r.URL.Path == "/api/directory/import" && f.refusals > 0 {
		f.refusals--
		http.NotFound(w, r)
		return
	}
	_, _ = io.WriteString(w, `{"version":1,"spaces":[]}`)
}

// A sync before the import would apply a snapshot that has never heard of
// Courrier's spaces, and remove their members. Run must keep retrying the
// import through Registre's 404s and pull a snapshot only afterwards.
func TestImportSucceedsBeforeTheFirstSync(t *testing.T) {
	fake := &fakeRegistre{refusals: 2}
	server := httptest.NewServer(fake)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	applied := make(chan struct{})
	syncer := &directory.Syncer{
		Client: directory.Client{Issuer: server.URL},
		Apply: func(context.Context, directory.Snapshot) error {
			close(applied)
			cancel()
			return nil
		},
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	load := func(context.Context) ([]directory.ImportSpace, error) { return nil, nil }
	importThenSync(ctx, syncer, load, time.Millisecond)

	select {
	case <-applied:
	default:
		t.Fatal("the snapshot was never applied")
	}
	fake.mu.Lock()
	defer fake.mu.Unlock()
	want := []string{"/api/directory/import", "/api/directory/import", "/api/directory/import", "/api/directory/snapshot"}
	if !slices.Equal(fake.calls, want) {
		t.Fatalf("calls = %v, want %v", fake.calls, want)
	}
}

// Registre matches an imported member by subject first, so a member who
// signed in through Registre must be sent by subject and everyone else by
// the email Registre falls back to.
func TestImportsNameMembersBySubjectElseEmail(t *testing.T) {
	dir, _ := directoryFixture(t)
	spaces, err := dir.imports(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(spaces) != 1 || spaces[0].Name != "Studio" || len(spaces[0].Members) != 4 {
		t.Fatalf("imports = %+v", spaces)
	}
	got := map[string]string{}
	for _, m := range spaces[0].Members {
		got[m.Subject+m.Email] = m.Role
	}
	want := map[string]string{"alice": "owner", "bob@facile.studio": "member", "carol@facile.studio": "member", "dave": "member"}
	for key, role := range want {
		if got[key] != role {
			t.Errorf("imports = %v, want %v", got, want)
		}
	}
}
