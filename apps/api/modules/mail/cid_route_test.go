package mail

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestCIDParamIsDecoded(t *testing.T) {
	const stored = "image001.png@01DD0333.B53E0920"

	var got string
	r := chi.NewRouter()
	r.Get("/cid/{cid}", func(w http.ResponseWriter, req *http.Request) {
		cid := chi.URLParam(req, "cid")
		if decoded, err := url.PathUnescape(cid); err == nil {
			cid = decoded
		}
		got = cid
	})

	req := httptest.NewRequest(http.MethodGet, "/cid/image001.png%4001DD0333.B53E0920", nil)
	r.ServeHTTP(httptest.NewRecorder(), req)

	if got != stored {
		t.Fatalf("cid = %q, want %q", got, stored)
	}
}
