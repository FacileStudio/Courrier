package mail

import (
	"context"
	"testing"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

func inboxSubjects(t *testing.T, env triageEnv, tab string) []string {
	t.Helper()
	category, err := tabCategory(tab)
	if err != nil {
		t.Fatal(err)
	}
	convs, total, err := env.svc.GetConversationsByFolderType(context.Background(), 1, env.a.ID,
		FolderQuery{FolderType: schemas.FolderTypeInbox, Page: 1, Limit: 50, Category: category})
	if err != nil {
		t.Fatal(err)
	}
	if int(total) != len(convs) {
		t.Fatalf("count %d disagrees with the %d rows listed", total, len(convs))
	}
	out := []string{}
	for _, c := range convs {
		out = append(out, c.Email.Subject)
	}
	return out
}

func TestInboxTabsAndSenderRules(t *testing.T) {
	env := newTriageEnv(t)
	testORM.Exec("DELETE FROM sender_rules")
	now := time.Now()
	addMail(t, mailKey{env.a.ID, env.inboxA, "t1"}, "lisa@epitech.eu", now, schemas.Email{Subject: "Visio"})
	addMail(t, mailKey{env.a.ID, env.inboxA, "t2"}, "news@railway.app", now.Add(-time.Hour), schemas.Email{Subject: "Railway news", Category: "bulk"})
	addMail(t, mailKey{env.a.ID, env.inboxA, "t3"}, "noreply@github.com", now.Add(-2*time.Hour), schemas.Email{Subject: "CI failed", Category: "auto"})

	if got := inboxSubjects(t, env, "primary"); len(got) != 1 || got[0] != "Visio" {
		t.Fatalf("primary = %v", got)
	}
	if got := inboxSubjects(t, env, "newsletters"); len(got) != 1 || got[0] != "Railway news" {
		t.Fatalf("newsletters = %v", got)
	}
	if got := inboxSubjects(t, env, ""); len(got) != 3 {
		t.Fatalf("no tab lists everything, got %v", got)
	}

	if err := env.svc.SetSenderRule(context.Background(), 1, env.a.ID, SenderRuleRequest{Address: "News@Railway.app", Tab: "primary"}); err != nil {
		t.Fatal(err)
	}
	if got := inboxSubjects(t, env, "primary"); len(got) != 2 {
		t.Fatalf("a primary pin must pull the newsletter into primary, got %v", got)
	}
	if err := env.svc.SetSenderRule(context.Background(), 1, env.a.ID, SenderRuleRequest{Address: "news@railway.app", Tab: "auto"}); err != nil {
		t.Fatal(err)
	}
	if got := inboxSubjects(t, env, "primary"); len(got) != 1 {
		t.Fatalf("auto must drop the pin, got %v", got)
	}
	if err := env.svc.SetSenderRule(context.Background(), 2, env.a.ID, SenderRuleRequest{Address: "x@y.z", Tab: "primary"}); err == nil {
		t.Fatal("another user's account must be refused")
	}
	if _, err := tabCategory("spam"); err == nil {
		t.Fatal("unknown tab must be refused")
	}
}

func TestSearchMatchesEveryWordAnywhere(t *testing.T) {
	env := newTriageEnv(t)
	addMail(t, mailKey{env.a.ID, env.inboxA, "t1"}, "lisa1.nay@epitech.eu", time.Now(), schemas.Email{Subject: "Visio de point", FromName: "NAY Lisa"})
	addMail(t, mailKey{env.a.ID, env.inboxA, "t2"}, "lisa@other.test", time.Now(), schemas.Email{Subject: "Hello", FromName: "Lisa Martin"})
	emails, _, err := env.svc.SearchEmails(context.Background(), 1, env.a.ID, "Lisa Nay", 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(emails) != 1 || emails[0].Subject != "Visio de point" {
		t.Fatalf("want only Lisa Nay's thread, got %d", len(emails))
	}
}
