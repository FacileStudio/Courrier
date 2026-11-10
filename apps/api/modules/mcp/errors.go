package mcp

import (
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/FacileStudio/tronc/errors"
)

// Hints appended to a not_found error, naming the tool that lists valid ids.
const (
	hintAccount = "Call courrier_list_accounts for the account ids you can use."
	hintMessage = "Call courrier_search_mail for the message ids of this account."
	hintDraft   = "Call courrier_search_mail with folder \"drafts\" for this account's draft ids, or create one with courrier_draft."
)

// errUnavailable is the text a model gets for a failure it cannot fix.
var errUnavailable = stderrors.New("Courrier could not finish this request. Retry once; if it fails again, tell the user")

// failure turns a service error into one a model can act on: the service's
// own message for a bad argument or a missing precondition, the message plus
// hint for a missing record, and a fixed text for anything internal, whose
// detail stays out of the result. A nil err stays nil.
func failure(err error, hint string) error {
	if err == nil {
		return nil
	}
	var known *errors.Error
	if !stderrors.As(err, &known) {
		return errUnavailable
	}
	switch known.Code {
	case "not_found":
		return fmt.Errorf("%s. %s", sentence(known.Message), hint)
	case "invalid_argument", "failed_precondition", "already_exists", "permission_denied":
		return fmt.Errorf("%s", sentence(known.Message))
	default:
		return errUnavailable
	}
}

// sentence capitalises the first letter of a service message.
func sentence(message string) string {
	if message == "" {
		return message
	}
	return strings.ToUpper(message[:1]) + message[1:]
}
