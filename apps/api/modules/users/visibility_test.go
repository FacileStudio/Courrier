package users

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
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/internal/middleware"
	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/kit"
	portepg "github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/porte/session"
	"github.com/FacileStudio/tronc/testdb"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

const bootstrap = "createdb courrier_test  # or point at any scratch database"

var testORM *gorm.DB

func TestMain(m *testing.M) {
	url, configured := testdb.URL()
	if !configured {
		testdb.Announce(bootstrap)
	} else {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_users_test", Migrate: schemas.Migrate})
		if err != nil {
			log.Fatalf("testdb: %v", err)
		}
		testORM = db
	}
	os.Exit(m.Run())
}

type tableResolver struct{ orm *gorm.DB }

func (r tableResolver) Identity(ctx context.Context, userID int64) (kit.Identity, error) {
	var user schemas.User
	if err := r.orm.WithContext(ctx).First(&user, userID).Error; err != nil {
		return kit.Identity{}, err
	}
	return kit.Identity{UserID: user.ID, Email: user.Email, Name: user.Name}, nil
}

type visibilityFixture struct {
	router                        chi.Router
	caller, peer, stranger, admin schemas.User
	callerToken, adminToken       string
}

func newVisibilityFixture(t *testing.T) visibilityFixture {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason(bootstrap))
	}
	for _, table := range []string{"porte_sessions", "space_members", "spaces", "users"} {
		if err := testORM.Exec("DELETE FROM " + table).Error; err != nil {
			t.Fatal(err)
		}
	}
	sqlDB, err := testORM.DB()
	if err != nil {
		t.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sessions, err := session.New(porte.Config{}, session.Deps{Sessions: portepg.New(sqlDB).Sessions(), Logger: logger})
	if err != nil {
		t.Fatal(err)
	}
	authService := auth.NewService(sessions, tableResolver{testORM})
	router := chi.NewRouter()
	RegisterRoutes(router, NewService(testORM, t.TempDir(), authService, middleware.Admins{"Admin@Example.test"}), authService)

	f := visibilityFixture{router: router,
		caller: seedUser(t, "caller@example.test"), peer: seedUser(t, "peer@example.test"),
		stranger: seedUser(t, "stranger@example.test"), admin: seedUser(t, "admin@example.test")}
	shareSpace(t, f.caller, f.peer)
	shareSpace(t, f.stranger)
	f.callerToken = issue(t, authService, f.caller.ID)
	f.adminToken = issue(t, authService, f.admin.ID)
	return f
}

func seedUser(t *testing.T, email string) schemas.User {
	t.Helper()
	user := schemas.User{Email: email, Name: email}
	if err := testORM.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	return user
}

func shareSpace(t *testing.T, members ...schemas.User) {
	t.Helper()
	space := schemas.Space{Name: "shared"}
	if err := testORM.Create(&space).Error; err != nil {
		t.Fatal(err)
	}
	for _, member := range members {
		if err := testORM.Create(&schemas.SpaceMember{SpaceID: space.ID, UserID: member.ID, Role: schemas.SpaceRoleMember}).Error; err != nil {
			t.Fatal(err)
		}
	}
}

func issue(t *testing.T, authService *auth.Service, userID int64) string {
	t.Helper()
	token, _, err := authService.Issue(context.Background(), userID, "test")
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func (f visibilityFixture) get(t *testing.T, token, path string) (int, []byte) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	recorder := httptest.NewRecorder()
	f.router.ServeHTTP(recorder, request)
	return recorder.Code, recorder.Body.Bytes()
}

func (f visibilityFixture) listed(t *testing.T, token string) map[string]bool {
	t.Helper()
	status, body := f.get(t, token, "/users/")
	if status != http.StatusOK {
		t.Fatalf("GET /users/ = %d %s", status, body)
	}
	var list ListResponse
	if err := json.Unmarshal(body, &list); err != nil {
		t.Fatal(err)
	}
	emails := map[string]bool{}
	for _, user := range list.Users {
		emails[user.Email] = true
	}
	return emails
}

// Any Registre account can sign in, so a session that is not on
// COURRIER_ADMINS must not read the address of someone it shares no space with.
func TestANonAdminSessionCannotListEveryUser(t *testing.T) {
	f := newVisibilityFixture(t)

	got := f.listed(t, f.callerToken)
	if len(got) != 2 || !got["caller@example.test"] || !got["peer@example.test"] {
		t.Fatalf("GET /users/ as a non-admin = %v, want only itself and its space peer", got)
	}

	for _, hidden := range []schemas.User{f.stranger, f.admin} {
		if status, _ := f.get(t, f.callerToken, "/users/"+strconv.FormatInt(hidden.ID, 10)); status != http.StatusNotFound {
			t.Errorf("GET /users/{%s} as a non-admin = %d, want 404", hidden.Email, status)
		}
	}
	if status, _ := f.get(t, f.callerToken, "/users/"+strconv.FormatInt(f.peer.ID, 10)); status != http.StatusOK {
		t.Errorf("GET /users/{peer} = %d, want 200", status)
	}
}

func TestAnUnauthenticatedRequestCannotListUsers(t *testing.T) {
	f := newVisibilityFixture(t)
	if status, body := f.get(t, "", "/users/"); status != http.StatusUnauthorized {
		t.Fatalf("GET /users/ without a token = %d %s, want 401", status, body)
	}
}

func TestAnAdminListsEveryUser(t *testing.T) {
	f := newVisibilityFixture(t)
	if got := f.listed(t, f.adminToken); len(got) != 4 {
		t.Fatalf("GET /users/ as an admin = %v, want all four accounts", got)
	}
	if status, _ := f.get(t, f.adminToken, "/users/"+strconv.FormatInt(f.stranger.ID, 10)); status != http.StatusOK {
		t.Fatalf("GET /users/{stranger} as an admin = %d, want 200", status)
	}
}
