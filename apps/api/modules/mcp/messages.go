package mcp

import (
	"fmt"
	"strings"

	"github.com/FacileStudio/Courrier/apps/api/modules/mail"
	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

// question is what the client asks the person.
func question(kind string, draft schemas.Email) string {
	recipients := strings.Join(append(append(mail.AddressList(draft.ToAddresses), mail.AddressList(draft.CcAddresses)...),
		mail.AddressList(draft.BccAddresses)...), ", ")
	if kind == kindSend {
		return fmt.Sprintf("Send the email %q to %s now?", draft.Subject, recipients)
	}
	return fmt.Sprintf("Delete the draft %q (to %s) for good?", draft.Subject, recipients)
}

// doneStatus is the status of a confirmed action that ran.
func doneStatus(kind string) string {
	if kind == kindSend {
		return "sent"
	}
	return "deleted"
}

// doneMessage tells the model what happened.
func doneMessage(kind string) string {
	if kind == kindSend {
		return "The user confirmed and the email was sent. A copy is in Sent."
	}
	return "The user confirmed and the draft was deleted."
}
