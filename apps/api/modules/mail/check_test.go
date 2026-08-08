package mail

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/FacileStudio/tronc/testdb"
	"gorm.io/gorm"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

var testORM *gorm.DB

func TestMain(m *testing.M) {
	url, configured := testdb.URL()
	if !configured {
		testdb.Announce("createdb courrier_test  # or point at any scratch database")
	} else {
		db, err := testdb.Open(url, testdb.Config{Prefix: "courrier_mail_test", Migrate: schemas.Migrate})
		if err != nil {
			log.Fatalf("testdb: %v", err)
		}
		testORM = db
	}
	os.Exit(m.Run())
}

/*
The invariant worth a test: both legs always run. TestConnection returns on the first
failure, and CheckAccount deliberately does not — "IMAP works, SMTP does not" is the answer
someone repairing an account needs, and a regression to fail-fast would silently halve it
while every other assertion here still passed.

127.0.0.1:1 is refused immediately, so this needs no network and no mail server.
*/
func TestCheckAccountReportsBothLegs(t *testing.T) {
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	if err := testORM.Exec(`DELETE FROM accounts`).Error; err != nil {
		t.Fatalf("clear accounts: %v", err)
	}

	service := NewService(testORM, nil)

	account := schemas.Account{
		UserID:       1,
		Name:         "Refused both ways",
		Email:        "someone@example.test",
		IMAPHost:     "127.0.0.1",
		IMAPPort:     1,
		IMAPUser:     "someone",
		IMAPPassword: "",
		SMTPHost:     "127.0.0.1",
		SMTPPort:     1,
		SMTPUser:     "someone",
		SMTPPassword: "",
	}
	if err := testORM.Create(&account).Error; err != nil {
		t.Fatalf("create account: %v", err)
	}

	result, err := service.CheckAccount(context.Background(), account.UserID, account.ID)
	if err != nil {
		t.Fatalf("CheckAccount returned a transport error, but a refused handshake is a result: %v", err)
	}

	for name, leg := range map[string]CheckLeg{"IMAP": result.IMAP, "SMTP": result.SMTP} {
		if !leg.Configured {
			t.Errorf("%s: Configured = false, want true (a host is set)", name)
		}
		if leg.OK {
			t.Errorf("%s: OK = true against a closed port", name)
		}
		if leg.Error == "" {
			t.Errorf("%s: no error text — this is the fail-fast regression, the second leg never ran", name)
		}
	}
}

// An unset host is a different answer from a failed one: the fix is a form field, not a network.
func TestCheckAccountLeavesUnsetProtocolUnconfigured(t *testing.T) {
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	if err := testORM.Exec(`DELETE FROM accounts`).Error; err != nil {
		t.Fatalf("clear accounts: %v", err)
	}

	service := NewService(testORM, nil)

	account := schemas.Account{
		UserID:   1,
		Name:     "Receive only",
		Email:    "someone@example.test",
		IMAPHost: "127.0.0.1",
		IMAPPort: 1,
	}
	if err := testORM.Create(&account).Error; err != nil {
		t.Fatalf("create account: %v", err)
	}

	result, err := service.CheckAccount(context.Background(), account.UserID, account.ID)
	if err != nil {
		t.Fatalf("CheckAccount: %v", err)
	}

	if !result.IMAP.Configured || result.IMAP.Error == "" {
		t.Errorf("IMAP leg = %+v, want configured with an error", result.IMAP)
	}
	if result.SMTP.Configured || result.SMTP.OK || result.SMTP.Error != "" {
		t.Errorf("SMTP leg = %+v, want the zero value for an unset host", result.SMTP)
	}
	if strings.Contains(result.IMAP.Error, "SMTP") {
		t.Errorf("IMAP error mentions SMTP: %q", result.IMAP.Error)
	}
}
