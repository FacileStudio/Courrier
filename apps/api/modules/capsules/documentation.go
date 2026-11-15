package capsules

import documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"

// Documentation describes the capsule relay for the API reference.
func Documentation() documentation.Module {
	return documentation.Module{
		Name: "capsules",
		Description: "Secrets sealed as Capsules that burn on their first read, for a link in a mail. The browser encrypts " +
			"and keeps the key in the link's fragment; Courrier relays only ciphertext to Capsule.",
		Routes: []documentation.Route{
			{
				Method:  "POST",
				Path:    "/capsules",
				Summary: "Seal a capsule",
				Description: "Stores the ciphertext in Capsule as a burn-after-read capsule. url is the Capsule link without " +
					"its key. Rate limited to 30 per minute.",
				Auth:         "session cookie or bearer token",
				RequestBody:  SealRequest{},
				ResponseBody: SealResponse{},
				Errors:       sealErrors(),
			},
		},
	}
}

func sealErrors() []documentation.Error {
	return []documentation.Error{
		{Status: 400, Code: "invalid_argument", Description: "Missing ciphertext or unknown expiry."},
		{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."},
		{Status: 413, Code: "resource_exhausted", Description: "Ciphertext over 1 MB."},
		{Status: 429, Code: "rate_limited", Description: "Too many capsules."},
		{Status: 503, Code: "unavailable", Description: "Capsule did not answer."},
	}
}
