package mail

import (
	"strings"
	"testing"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func TestWithoutBccStripsTheHeaderAndItsFolds(t *testing.T) {
	spec := messageSpec{
		FromEmail: "me@example.test", To: []string{"to@x.test"}, Subject: "s", Body: "hello",
		Bcc: []string{"Hidden One <h1@x.test>", "h2@x.test", "h3@x.test", "h4@x.test", "h5@x.test"},
	}
	raw, _, err := buildMessage(spec)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "Bcc:") {
		t.Fatalf("the filed copy must carry Bcc:\n%s", raw)
	}
	clean := string(withoutBcc(raw))
	for _, leak := range []string{"Bcc:", "h1@x.test", "h5@x.test"} {
		if strings.Contains(clean, leak) {
			t.Errorf("SMTP copy still contains %q:\n%s", leak, clean)
		}
	}
	for _, keep := range []string{"To: <to@x.test>", "Subject: s", "hello"} {
		if !strings.Contains(clean, keep) {
			t.Errorf("SMTP copy lost %q:\n%s", keep, clean)
		}
	}
}

func TestBuildMessageAllowsBccOnly(t *testing.T) {
	raw, _, err := buildMessage(messageSpec{FromEmail: "me@example.test", Bcc: []string{"h@x.test"}, Body: "b"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "\r\nTo:") {
		t.Errorf("an empty To must not be written:\n%s", raw)
	}
}

func TestEnvelopeRecipientsIncludeBccOnceEach(t *testing.T) {
	req := SendRequest{To: []string{"A <a@x.test>"}, Cc: []string{"a@X.test", "c@x.test"}, Bcc: []string{"b@x.test", "C@x.test"}}
	got := strings.Join(envelopeRecipients(req), ",")
	if got != "a@x.test,c@x.test,b@x.test" {
		t.Errorf("envelope = %s", got)
	}
}

func TestSendHidesBccFromTheDeliveredCopyAndRecordsIt(t *testing.T) {
	env := newMailEnv(t)
	req := SendRequest{To: []string{"to@x.test"}, Bcc: []string{"hidden@x.test"}, Subject: "blind", Body: "hi"}
	if err := env.svc.Send(ctx(), 1, env.account.ID, req); err != nil {
		t.Fatal(err)
	}
	if len(*env.sent) != 1 {
		t.Fatalf("delivered %d messages", len(*env.sent))
	}
	got := (*env.sent)[0]
	if strings.Join(got.recipients, ",") != "to@x.test,hidden@x.test" {
		t.Errorf("RCPT = %v", got.recipients)
	}
	if strings.Contains(string(got.msg), "hidden@x.test") {
		t.Errorf("the delivered copy names the blind recipient:\n%s", got.msg)
	}
	var row schemas.Email
	if err := testORM.Where("subject = ?", "blind").First(&row).Error; err != nil {
		t.Fatal(err)
	}
	resp := emailToResponse(row)
	if len(resp.BccAddresses) != 1 || resp.BccAddresses[0].Email != "hidden@x.test" {
		t.Errorf("bcc_addresses = %+v", resp.BccAddresses)
	}
	waitFor(t, "the Sent copy", func() bool { return env.fake.count(t, "Sent") == 1 })
}

func TestSendRequiresSomeRecipient(t *testing.T) {
	env := newMailEnv(t)
	if err := env.svc.Send(ctx(), 1, env.account.ID, SendRequest{Body: "x"}); err == nil {
		t.Fatal("a message with no To, Cc or Bcc must be refused")
	}
	if err := env.svc.Send(ctx(), 1, env.account.ID, SendRequest{Bcc: []string{"b@x.test"}, Body: "x"}); err != nil {
		t.Fatalf("a Bcc-only message is valid: %v", err)
	}
}

func TestNullBccColumnReadsAsEmptyList(t *testing.T) {
	env := newMailEnv(t)
	email := env.seed(t, "INBOX", "legacy row")
	if err := testORM.Exec("UPDATE emails SET bcc_addresses = NULL WHERE id = ?", email.ID).Error; err != nil {
		t.Fatal(err)
	}
	got, err := env.svc.GetEmail(ctx(), 1, env.account.ID, email.ID)
	if err != nil {
		t.Fatal(err)
	}
	if resp := emailToResponse(got); resp.BccAddresses == nil || len(resp.BccAddresses) != 0 {
		t.Errorf("bcc_addresses = %#v, want []", resp.BccAddresses)
	}
}
