package mcp

// toolOf is the tool that asks for a pending action kind.
func toolOf(kind string) string {
	if kind == kindSend {
		return toolSend
	}
	return toolDeleteDraft
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

// waiting tells the model that nothing happened yet, and what to do.
func waiting(kind string) string {
	what := "Nothing was sent"
	if kind == kindDeleteDraft {
		what = "Nothing was deleted"
	}
	return what + ": your client cannot ask the user to confirm. Give the user confirm_url and ask them to " +
		"confirm there. Do not call this tool again for the same draft; the request stays open for 7 days."
}
