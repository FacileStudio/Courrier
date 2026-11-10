package mcp

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/testdb"
	"gorm.io/gorm"
)

// subject is the Registre sub of the test person.
const subject = "person-sub"

// testORM is the scratch database, nil when TEST_DATABASE_URL is unset.
var testORM *gorm.DB

// TestMain opens the scratch database when one is configured.
func TestMain(m *testing.M) {
	if url, configured := testdb.URL(); configured {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_mcp_test", Migrate: schemas.Migrate})
		if err != nil {
			log.Fatalf("testdb: %v", err)
		}
		testORM = db
	}
	os.Exit(m.Run())
}

// needDB skips a test that stores pending actions when there is no database.
func needDB(t *testing.T) {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
}

// registre is a fake authorization server: discovery, JWKS and a minter.
type registre struct {
	srv *httptest.Server
	key *rsa.PrivateKey
}

// newRegistre serves discovery and one RSA key.
func newRegistre(t *testing.T) *registre {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	r := &registre{key: key}
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"issuer": r.srv.URL, "jwks_uri": r.srv.URL + "/jwks"})
	})
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{publicJWK(&key.PublicKey)}})
	})
	r.srv = httptest.NewServer(mux)
	t.Cleanup(r.srv.Close)
	return r
}

// publicJWK is the JWKS entry of key.
func publicJWK(key *rsa.PublicKey) map[string]string {
	return map[string]string{
		"kty": "RSA", "kid": "k1", "n": base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
	}
}

// mint signs an access token for resource; overrides replace claims, and a
// nil value removes one.
func (r *registre) mint(t *testing.T, resource string, overrides map[string]any) string {
	t.Helper()
	claims := map[string]any{
		"iss": r.srv.URL, "sub": subject, "aud": []string{resource}, "client_id": "https://claude.ai/oauth/mcp.json",
		"exp": time.Now().Add(10 * time.Minute).Unix(), "iat": time.Now().Unix(), "jti": "j1",
		"scope": "courrier.read courrier.write",
	}
	for name, value := range overrides {
		claims[name] = value
		if value == nil {
			delete(claims, name)
		}
	}
	header, _ := json.Marshal(map[string]string{"alg": "RS256", "typ": "at+jwt", "kid": "k1"})
	body, _ := json.Marshal(claims)
	input := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(body)
	sum := sha256.Sum256([]byte(input))
	signature, err := rsa.SignPKCS1v15(rand.Reader, r.key, crypto.SHA256, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	return input + "." + base64.RawURLEncoding.EncodeToString(signature)
}

// identities maps the test subject to userID and nothing else.
type identities struct {
	porte.IdentityStore
	userID int64
}

// Find resolves only the test subject.
func (s identities) Find(_ context.Context, provider, sub string) (porte.StoredIdentity, error) {
	if sub != subject {
		return porte.StoredIdentity{}, porte.ErrNotFound
	}
	return porte.StoredIdentity{UserID: s.userID, Provider: provider, Subject: sub}, nil
}

// seedPerson creates the user and account rows pending actions point at,
// or returns fixed ids without a database.
func seedPerson(t *testing.T) (userID, accountID int64) {
	t.Helper()
	if testORM == nil {
		return 7, 70
	}
	user := schemas.User{Email: t.Name() + "-" + rand.Text() + "@example.test"}
	if err := testORM.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	account := schemas.Account{UserID: user.ID, Email: "me@example.test", Name: "Me"}
	if err := testORM.Create(&account).Error; err != nil {
		t.Fatal(err)
	}
	return user.ID, account.ID
}
