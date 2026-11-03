package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FacileStudio/porte/directory"
)

// Registre posts its nudge with no session, no token and no Origin. Anything
// on the way that wants one turns every change in Registre into a five-minute
// wait for the next scheduled pull, silently.
func TestRegistreNudgeNeedsNoCredentials(t *testing.T) {
	router := referenceRouter(t)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, directory.SyncPath, nil))

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("POST %s = %d, want 202", directory.SyncPath, recorder.Code)
	}
}
