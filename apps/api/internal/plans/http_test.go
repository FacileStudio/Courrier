package plans

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/porte/forfait"
	"github.com/FacileStudio/tronc/errors"
)

// TestWriteErrorShape pins the 402 envelope and that an unsafe upgrade link
// never reaches the client.
func TestWriteErrorShape(t *testing.T) {
	cases := map[string]string{
		"https://sso.test/spaces/spc_1": "https://sso.test/spaces/spc_1",
		"http://sso.test/spaces":        "",
		"javascript:alert(1)":           "",
		"https://u:p@sso.test/spaces":   "",
		"":                              "",
	}
	for link, want := range cases {
		recorder := httptest.NewRecorder()
		WriteError(recorder, &forfait.LimitError{Key: Mailboxes, Limit: 1, UpgradeURL: link})
		var body struct {
			Error map[string]any `json:"error"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		got, _ := body.Error["upgrade_url"].(string)
		if recorder.Code != http.StatusPaymentRequired || body.Error["code"] != "payment_required" ||
			body.Error["key"] != Mailboxes || body.Error["limit"] != float64(1) || got != want {
			t.Errorf("%q: status %d, body %v", link, recorder.Code, body.Error)
		}
	}
	recorder := httptest.NewRecorder()
	WriteError(recorder, errors.Forbidden("no"))
	if recorder.Code != http.StatusForbidden {
		t.Errorf("other errors: status %d", recorder.Code)
	}
}
