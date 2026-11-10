package ideas

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/internal/plans"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte/kit"

	"github.com/FacileStudio/porte/directory"
	"github.com/FacileStudio/tronc/bulle"
	"github.com/FacileStudio/tronc/testdb"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

const testIssuer = "https://sso.test"

var testORM *gorm.DB

func TestMain(m *testing.M) {
	if url, ok := testdb.URL(); ok {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_ideas_test", Migrate: schemas.Migrate})
		if err != nil {
			log.Fatalf("testdb: %v", err)
		}
		testORM = db
	}
	os.Exit(m.Run())
}

type fakeBulle struct {
	mu    sync.Mutex
	calls []string
	space []string
}

func (f *fakeBulle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	f.calls = append(f.calls, r.Method+" "+r.URL.Path)
	f.space = append(f.space, r.URL.Query().Get("space_id"))
	f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"ideas":[],"generating":false}`))
}

func (f *fakeBulle) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

type env struct {
	router chi.Router
	bulle  *fakeBulle
	gate   *plans.Gate
}

func setup(t *testing.T) env {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test"))
	}
	for _, table := range []string{"forfait_state", "accounts", "space_members", "spaces", "porte_identities", "users"} {
		if err := testORM.Exec("DELETE FROM " + table).Error; err != nil {
			t.Fatal(err)
		}
	}
	fake := &fakeBulle{}
	server := httptest.NewServer(fake)
	t.Cleanup(server.Close)
	gate := plans.New(testORM, testIssuer, slog.New(slog.NewTextHandler(io.Discard, nil)))
	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					id, err := strconv.ParseInt(req.Header.Get("X-Test-User"), 10, 64)
					if err != nil {
						http.Error(w, "unauthenticated", http.StatusUnauthorized)
						return
					}
					next.ServeHTTP(w, req.WithContext(kit.WithIdentity(req.Context(), kit.Identity{UserID: id})))
				})
			})
			err := New(testORM, gate).Mount(r, bulle.Options{
				URL:         server.URL,
				TokenSource: func(context.Context) (string, error) { return "tok", nil },
				UpgradeURL:  testIssuer + "/spaces",
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	})
	return env{router: router, bulle: fake, gate: gate}
}

func (e env) do(method, target, user string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, nil)
	if user != "" {
		req.Header.Set("X-Test-User", user)
	}
	rec := httptest.NewRecorder()
	e.router.ServeHTTP(rec, req)
	return rec
}

func newUser(t *testing.T, email, subject string) string {
	t.Helper()
	row := schemas.User{Email: email}
	if err := testORM.Create(&row).Error; err != nil {
		t.Fatal(err)
	}
	if err := testORM.Exec(`INSERT INTO porte_identities (user_id, provider, subject) VALUES (?, ?, ?)`,
		row.ID, testIssuer, subject).Error; err != nil {
		t.Fatal(err)
	}
	return strconv.FormatInt(row.ID, 10)
}

func newSpace(t *testing.T, link string, members ...string) string {
	t.Helper()
	row := schemas.Space{Name: "S"}
	if link != "" {
		row.RegistreSpaceID = &link
	}
	if err := testORM.Create(&row).Error; err != nil {
		t.Fatal(err)
	}
	for _, id := range members {
		if err := testORM.Exec(`INSERT INTO space_members (space_id, user_id, role) VALUES (?, ?, 'member')`, row.ID, id).Error; err != nil {
			t.Fatal(err)
		}
	}
	return row.ID
}

func entitle(t *testing.T, gate *plans.Gate, spaces, people []string) {
	t.Helper()
	on := map[string]json.RawMessage{plans.AIAccess: json.RawMessage("true")}
	snap := directory.Snapshot{Version: 100, EntitlementsRevision: 1}
	for _, id := range spaces {
		snap.Spaces = append(snap.Spaces, directory.Space{ID: id, Entitlements: on})
	}
	for _, sub := range people {
		snap.People = append(snap.People, directory.Person{Subject: sub, Entitlements: on})
	}
	if len(snap.Spaces)+len(snap.People) == 0 {
		snap.People = append(snap.People, directory.Person{Subject: "nobody"})
	}
	if err := gate.Apply(context.Background(), snap); err != nil {
		t.Fatal(err)
	}
}

func TestRequiresAuth(t *testing.T) {
	e := setup(t)
	if rec := e.do("GET", "/api/ideas", ""); rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
	if e.bulle.count() != 0 {
		t.Fatal("Bulle was called without a session")
	}
}

func TestFreePlanAnswers402(t *testing.T) {
	e := setup(t)
	alice := newUser(t, "alice@facile.studio", "alice")
	entitle(t, e.gate, nil, nil)
	list := e.do("GET", "/api/ideas", alice)
	var page struct {
		AIAccess bool `json:"ai_access"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &page); err != nil || page.AIAccess {
		t.Fatalf("list = %d %s, want ai_access false", list.Code, list.Body.String())
	}
	if rec := e.do("POST", "/api/ideas/generate", alice); rec.Code != http.StatusPaymentRequired {
		t.Fatalf("generate = %d, want 402", rec.Code)
	}
	if rec := e.do("POST", "/api/ideas/abc/run", alice); rec.Code != http.StatusPaymentRequired {
		t.Fatalf("run = %d, want 402", rec.Code)
	}
	if e.bulle.count() != 0 {
		t.Fatalf("Bulle was called for a free plan: %v", e.bulle.calls)
	}
}

func TestPersonalSendsNoSpace(t *testing.T) {
	e := setup(t)
	alice := newUser(t, "alice@facile.studio", "alice")
	entitle(t, e.gate, nil, []string{"alice"})
	if rec := e.do("POST", "/api/ideas/generate", alice); rec.Code == http.StatusPaymentRequired || rec.Code >= 500 {
		t.Fatalf("generate = %d %s", rec.Code, rec.Body.String())
	}
	if e.bulle.count() != 1 || e.bulle.space[0] != "" {
		t.Fatalf("calls = %v spaces = %q, want one call without a space", e.bulle.calls, e.bulle.space)
	}
}

func TestSpaceIsPassedAsRegistreID(t *testing.T) {
	e := setup(t)
	alice := newUser(t, "alice@facile.studio", "alice")
	bob := newUser(t, "bob@facile.studio", "bob")
	linked := newSpace(t, "spc_pro", alice)
	unlinked := newSpace(t, "", alice)
	entitle(t, e.gate, []string{"spc_pro"}, nil)
	if rec := e.do("GET", "/api/ideas?space_id="+linked, alice); rec.Code != http.StatusOK {
		t.Fatalf("member = %d %s", rec.Code, rec.Body.String())
	}
	if e.bulle.space[0] != "spc_pro" {
		t.Fatalf("Bulle saw space %q, want spc_pro", e.bulle.space[0])
	}
	cases := map[string]struct {
		target, user string
		want         int
	}{
		"non member":   {"/api/ideas?space_id=" + linked, bob, http.StatusNotFound},
		"unlinked":     {"/api/ideas?space_id=" + unlinked, alice, http.StatusPreconditionFailed},
		"bad spelling": {"/api/ideas?space_id=%7B" + linked + "%7D", alice, http.StatusBadRequest},
	}
	before := e.bulle.count()
	for name, c := range cases {
		if rec := e.do("GET", c.target, c.user); rec.Code != c.want {
			t.Errorf("%s = %d, want %d", name, rec.Code, c.want)
		}
	}
	if e.bulle.count() != before {
		t.Fatal("Bulle was called for a refused space")
	}
}

func TestSpacePlanDecidesOverPerson(t *testing.T) {
	e := setup(t)
	alice := newUser(t, "alice@facile.studio", "alice")
	free := newSpace(t, "spc_free", alice)
	entitle(t, e.gate, nil, []string{"alice"})
	if rec := e.do("POST", "/api/ideas/generate?space_id="+free, alice); rec.Code != http.StatusPaymentRequired {
		t.Fatalf("generate in a free space = %d, want 402", rec.Code)
	}
}
