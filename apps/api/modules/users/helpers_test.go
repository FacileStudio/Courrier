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
	"strings"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/modules/auth"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/porte/local"
	portepg "github.com/FacileStudio/porte/pg"
	"github.com/FacileStudio/porte/session"
	"github.com/FacileStudio/tronc/testdb"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

const testBootstrap = "createdb courrier_test  # or point at any scratch database"

var testORM *gorm.DB

func TestMain(m *testing.M) {
	url, configured := testdb.URL()
	if !configured {
		testdb.Announce(testBootstrap)
	} else {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_test", Migrate: schemas.Migrate})
		if err != nil {
			log.Fatalf("testdb: %v", err)
		}
		testORM = db
	}
	os.Exit(m.Run())
}

// harness is the users module behind the router that actually serves it.
//
// porte's ChangePassword reads the caller's session id out of the request
// context and writes the rotated cookie through the ResponseWriter, so a test
// calling the service with a bare context exercises neither and passes
// regardless.
type harness struct {
	t        *testing.T
	router   chi.Router
	sessions *session.Manager
	orm      *gorm.DB
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason(testBootstrap))
	}
	if err := testORM.Exec(`DELETE FROM users`).Error; err != nil {
		t.Fatalf("clear users: %v", err)
	}

	handle, err := testORM.DB()
	if err != nil {
		t.Fatalf("sql handle: %v", err)
	}
	store := portepg.New(handle)
	accounts := auth.NewUserStore(testORM)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	sessions, err := session.New(porte.Config{}, session.Deps{Sessions: store.Sessions(), Logger: logger})
	if err != nil {
		t.Fatalf("session manager: %v", err)
	}
	passwords, err := local.New(local.Config{AllowRegistration: true, MinPasswordLength: 8}, local.Deps{
		Users:      accounts,
		Identities: store.Identities(),
		Sessions:   sessions,
		Logger:     logger,
		Count:      accounts.CountUsers,
	})
	if err != nil {
		t.Fatalf("local kit: %v", err)
	}

	authService := auth.NewService(testORM, sessions, passwords, logger)
	router := chi.NewRouter()
	passwords.Mount(router)
	RegisterRoutes(router, NewService(testORM, t.TempDir(), authService), authService)

	return &harness{t: t, router: router, sessions: sessions, orm: testORM}
}

// send drives the real router, presenting the session as the browser does:
// a cookie plus the CSRF header porte requires of every mutating request that
// carries one.
func (h *harness) send(method, path, body, token string) *httptest.ResponseRecorder {
	h.t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	if token != "" {
		request.AddCookie(&http.Cookie{Name: porte.SessionCookieName, Value: token})
		request.Header.Set(porte.CSRFHeaderName, "1")
	}
	recorder := httptest.NewRecorder()
	h.router.ServeHTTP(recorder, request)
	return recorder
}

// register creates an account through porte's own registration route, which is
// what gives it the id-keyed local identity the migration is about.
func (h *harness) register(email, password string) (int64, string) {
	h.t.Helper()
	recorder := h.send(http.MethodPost, porte.RouteRegister, `{"email":"`+email+`","password":"`+password+`"}`, "")
	if recorder.Code != http.StatusCreated {
		h.t.Fatalf("register %s: %d %s", email, recorder.Code, recorder.Body.String())
	}
	var response porte.ExchangeResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		h.t.Fatalf("decode register response: %v", err)
	}
	var record schemas.User
	if err := h.orm.Where("email = ?", email).First(&record).Error; err != nil {
		h.t.Fatalf("read the registered account: %v", err)
	}
	return record.ID, response.Token
}

// federated creates an account with no password at all, which is the only
// shape SetPassword may write to.
func (h *harness) federated(email string) (int64, string) {
	h.t.Helper()
	record := schemas.User{Email: email, Name: "Federated"}
	if err := h.orm.Create(&record).Error; err != nil {
		h.t.Fatalf("create the federated account: %v", err)
	}
	return record.ID, h.login(record.ID)
}

// login issues an unlabelled session, which is what a second browser holds.
func (h *harness) login(userID int64) string {
	h.t.Helper()
	token, _, err := h.sessions.Issue(context.Background(), userID, "")
	if err != nil {
		h.t.Fatalf("issue a session: %v", err)
	}
	return token
}

// rotatedToken is the session porte minted while answering, read off the
// Set-Cookie header the way the browser reads it.
func (h *harness) rotatedToken(recorder *httptest.ResponseRecorder) string {
	h.t.Helper()
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == porte.SessionCookieName && cookie.Value != "" {
			return cookie.Value
		}
	}
	h.t.Fatalf("no rotated session cookie in the response")
	return ""
}

func errorCode(t *testing.T, recorder *httptest.ResponseRecorder) string {
	t.Helper()
	var envelope struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode the error envelope: %v", err)
	}
	return envelope.Error.Code
}

func errorMessage(t *testing.T, recorder *httptest.ResponseRecorder) string {
	t.Helper()
	var envelope struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode the error envelope: %v", err)
	}
	return envelope.Error.Message
}
