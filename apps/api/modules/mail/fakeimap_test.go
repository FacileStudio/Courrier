package mail

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/FacileStudio/tronc/testdb"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/emersion/go-imap/v2/imapserver/imapmemserver"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

type fakeIMAP struct {
	addr string
	user *imapmemserver.User
}

type capturedMail struct {
	recipients []string
	msg        []byte
}

type mailEnv struct {
	svc     *Service
	account schemas.Account
	fake    *fakeIMAP
	sent    *[]capturedMail
}

func startFakeIMAP(t *testing.T) *fakeIMAP {
	t.Helper()
	mem := imapmemserver.New()
	user := imapmemserver.NewUser("me", "pw")
	mem.AddUser(user)
	srv := imapserver.New(&imapserver.Options{
		NewSession: func(*imapserver.Conn) (imapserver.Session, *imapserver.GreetingData, error) {
			return mem.NewSession(), nil, nil
		},
		Caps:         imap.CapSet{imap.CapIMAP4rev1: {}, imap.CapUIDPlus: {}, imap.CapMove: {}, imap.CapChildren: {}},
		InsecureAuth: true,
	})
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go srv.Serve(ln)
	t.Cleanup(func() { srv.Close() })
	return &fakeIMAP{addr: ln.Addr().String(), user: user}
}

func (f *fakeIMAP) dial(string, int, string, string) (*imapclient.Client, error) {
	client, err := imapclient.DialInsecure(f.addr, nil)
	if err != nil {
		return nil, err
	}
	if err := client.Login("me", "pw").Wait(); err != nil {
		return nil, err
	}
	return client, nil
}

func (f *fakeIMAP) client(t *testing.T) *imapclient.Client {
	t.Helper()
	client, err := f.dial("", 0, "", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { closeIMAP(client) })
	return client
}

func (f *fakeIMAP) put(t *testing.T, mailbox, raw string) {
	t.Helper()
	client := f.client(t)
	cmd := client.Append(mailbox, int64(len(raw)), nil)
	cmd.Write([]byte(raw))
	cmd.Close()
	if _, err := cmd.Wait(); err != nil {
		t.Fatalf("append to %s: %v", mailbox, err)
	}
}

func (f *fakeIMAP) count(t *testing.T, mailbox string) int {
	t.Helper()
	status, err := f.client(t).Status(mailbox, &imap.StatusOptions{NumMessages: true}).Wait()
	if err != nil {
		t.Fatalf("status %s: %v", mailbox, err)
	}
	return int(*status.NumMessages)
}

func (f *fakeIMAP) mailboxes(t *testing.T) []string {
	t.Helper()
	boxes, err := f.client(t).List("", "*", nil).Collect()
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, box := range boxes {
		names = append(names, box.Mailbox)
	}
	return names
}

func (f *fakeIMAP) uids(t *testing.T, mailbox string) []imap.UID {
	t.Helper()
	client := f.client(t)
	if _, err := client.Select(mailbox, nil).Wait(); err != nil {
		t.Fatal(err)
	}
	found, err := client.UIDSearch(&imap.SearchCriteria{}, nil).Wait()
	if err != nil {
		t.Fatal(err)
	}
	return found.AllUIDs()
}

func ctx() context.Context { return context.Background() }

func waitFor(t *testing.T, what string, ok func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if ok() {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

func newMailEnv(t *testing.T) *mailEnv {
	t.Helper()
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	for _, table := range []string{"attachments", "emails", "thread_links", "contacts", "folders", "accounts"} {
		if err := testORM.Exec("DELETE FROM " + table).Error; err != nil {
			t.Fatalf("clear %s: %v", table, err)
		}
	}
	fake := startFakeIMAP(t)
	for _, name := range []string{"INBOX", "Drafts", "Sent", "Trash", "Archive", "Projects"} {
		if err := fake.user.Create(name, nil); err != nil && name != "INBOX" {
			t.Fatal(err)
		}
	}
	account := schemas.Account{UserID: 1, Name: "Me", Email: "me@example.test"}
	if err := testORM.Create(&account).Error; err != nil {
		t.Fatal(err)
	}
	env := &mailEnv{svc: NewService(testORM, nil), account: account, fake: fake, sent: &[]capturedMail{}}
	env.svc.dial = fake.dial
	env.svc.deliver = func(_ schemas.Account, rcpt []string, msg []byte) error {
		*env.sent = append(*env.sent, capturedMail{recipients: rcpt, msg: msg})
		return nil
	}
	if err := env.svc.SyncAccount(ctx(), 1, account.ID); err != nil {
		t.Fatalf("sync: %v", err)
	}
	return env
}

func (e *mailEnv) folder(t *testing.T, name string) schemas.Folder {
	t.Helper()
	var folder schemas.Folder
	if err := testORM.Where("account_id = ? AND path = ?", e.account.ID, name).First(&folder).Error; err != nil {
		t.Fatalf("folder %s: %v", name, err)
	}
	return folder
}

func (e *mailEnv) seed(t *testing.T, mailbox, subject string) schemas.Email {
	t.Helper()
	raw := fmt.Sprintf("From: a@x.test\r\nTo: me@example.test\r\nSubject: %s\r\nMessage-ID: <%s@x.test>\r\n\r\nbody of %s\r\n",
		subject, strings.ReplaceAll(subject, " ", "-"), subject)
	e.fake.put(t, mailbox, raw)
	folder := e.folder(t, mailbox)
	if err := e.svc.SyncFolderEmails(ctx(), 1, e.account.ID, folder.ID); err != nil {
		t.Fatal(err)
	}
	var email schemas.Email
	if err := testORM.Where("folder_id = ? AND subject = ?", folder.ID, subject).First(&email).Error; err != nil {
		t.Fatalf("seeded email %q: %v", subject, err)
	}
	return email
}

func connectIMAPRefused(string, int, string, string) (*imapclient.Client, error) {
	return nil, fmt.Errorf("connection refused")
}

func (f *fakeIMAP) flagged(t *testing.T, mailbox string) int {
	t.Helper()
	client := f.client(t)
	if _, err := client.Select(mailbox, nil).Wait(); err != nil {
		t.Fatal(err)
	}
	found, err := client.UIDSearch(&imap.SearchCriteria{Flag: []imap.Flag{imap.FlagFlagged}}, nil).Wait()
	if err != nil {
		t.Fatal(err)
	}
	return len(found.AllUIDs())
}
