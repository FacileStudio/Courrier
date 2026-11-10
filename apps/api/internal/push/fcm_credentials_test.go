package push

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"testing"
)

func serviceAccountJSON(t *testing.T, kind string) []byte {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	block := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: mustPKCS8(t, key)})
	raw, err := json.Marshal(map[string]string{
		"type": kind, "project_id": "synthetic-project", "private_key_id": "k1", "private_key": string(block),
		"client_email": "sender@synthetic-project.iam.gserviceaccount.com", "token_uri": "https://oauth2.example.test/token",
	})
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func mustPKCS8(t *testing.T, key *rsa.PrivateKey) []byte {
	t.Helper()
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return der
}

func TestNewFCMReadsAServiceAccountAndRefusesOtherKinds(t *testing.T) {
	sender, err := NewFCM(context.Background(), serviceAccountJSON(t, "service_account"))
	if err != nil {
		t.Fatalf("service account: %v", err)
	}
	if sender.endpoint != "https://fcm.googleapis.com/v1/projects/synthetic-project/messages:send" {
		t.Fatalf("endpoint = %q", sender.endpoint)
	}
	if _, err := NewFCM(context.Background(), serviceAccountJSON(t, "authorized_user")); err == nil {
		t.Fatal("a credential that is not a service account was accepted")
	}
}
