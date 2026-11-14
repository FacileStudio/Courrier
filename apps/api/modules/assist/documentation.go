package assist

import documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"

var unauthenticated = documentation.Error{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."}

// Documentation describes the writing helper routes for the API reference.
var Documentation = documentation.Module{
	Name: "assist",
	Description: "The Bulle writing helper. Courrier trades the user's Registre session for a Bulle token and calls " +
		"Bulle's /api/ai/complete as the user; Courrier holds no model key. It needs a personal plan with ai.access. " +
		"Off until BULLE_URL and OIDC are set.",
	Routes: []documentation.Route{
		{
			Method:  "GET",
			Path:    "/assist",
			Summary: "Say whether the helper can run",
			Description: "available is false with reason off (not configured), plan_required (the plan has no ai.access; " +
				"upgrade_url says where to get it), not_linked (no Registre session or no exchange grant) or unreachable.",
			Auth:         "session cookie or bearer token",
			ResponseBody: Status{},
			Errors:       []documentation.Error{unauthenticated},
		},
		{
			Method:  "POST",
			Path:    "/assist",
			Summary: "Run a helper action",
			Description: "action is rewrite, shorten, expand, fix, formal, friendly, translate (needs language), " +
				"reply or summarize (both need the thread as context). text is the draft or the selection. " +
				"Returns the new text. Rate limited to 30 per minute.",
			Auth:         "session cookie or bearer token",
			RequestBody:  Request{},
			ResponseBody: Response{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Unknown action, empty text, or missing language or context."},
				unauthenticated,
				{Status: 402, Code: "payment_required", Description: "The plan has no ai.access; upgrade_url says where to get it."},
				{Status: 412, Code: "failed_precondition", Description: "The helper is off, or the account is not linked."},
				{Status: 413, Code: "resource_exhausted", Description: "Text or context too long."},
				{Status: 429, Code: "rate_limited", Description: "Too many helper calls, or Bulle's daily limit is reached."},
				{Status: 503, Code: "unavailable", Description: "Bulle did not answer."},
			},
		},
	},
}
