package mail

import (
	"bufio"
	"bytes"
	"net/textproto"
	"regexp"
	"strings"

	"github.com/emersion/go-imap/v2/imapclient"
)

// Signals holds the header values that mark a message as bulk or automated.
type Signals struct {
	ListUnsubscribe string
	ListID          string
	Precedence      string
	AutoSubmitted   string
}

var noReplyLocal = regexp.MustCompile(`^(no-?reply|do-?not-?reply|notifications?|notify|mailer-daemon|postmaster|newsletter|news|marketing|updates)$`)

func bulkSignals(msg *imapclient.FetchMessageBuffer) Signals {
	raw := msg.FindBodySection(referencesSection)
	if len(raw) == 0 {
		return Signals{}
	}
	hdr, _ := textproto.NewReader(bufio.NewReader(bytes.NewReader(raw))).ReadMIMEHeader()
	return Signals{
		ListUnsubscribe: strings.TrimSpace(hdr.Get("List-Unsubscribe")),
		ListID:          strings.TrimSpace(hdr.Get("List-Id")),
		Precedence:      strings.ToLower(strings.TrimSpace(hdr.Get("Precedence"))),
		AutoSubmitted:   strings.ToLower(strings.TrimSpace(hdr.Get("Auto-Submitted"))),
	}
}

// Classify returns the category ("", "bulk" or "auto") and the reason that
// produced it. Header signals win over the sender pattern.
func Classify(s Signals, fromAddress string) (string, string) {
	switch {
	case s.ListUnsubscribe != "":
		return "bulk", "list-unsubscribe"
	case s.ListID != "":
		return "bulk", "list-id"
	case s.Precedence == "bulk" || s.Precedence == "list" || s.Precedence == "junk":
		return "bulk", "precedence-bulk"
	case s.AutoSubmitted != "" && s.AutoSubmitted != "no":
		return "auto", "auto-submitted"
	case isNoReplySender(fromAddress):
		return "bulk", "noreply-sender"
	}
	return "", ""
}

func isNoReplySender(address string) bool {
	local, _, ok := strings.Cut(strings.ToLower(strings.TrimSpace(address)), "@")
	return ok && noReplyLocal.MatchString(local)
}
