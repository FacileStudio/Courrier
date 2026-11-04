package plans

import (
	stderrors "errors"
	"net/http"
	"net/url"

	"github.com/FacileStudio/porte/forfait"
	"github.com/FacileStudio/tronc/httpjson"
)

// limitBody is the suite error envelope's inner object for a plan refusal.
type limitBody struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	UpgradeURL string `json:"upgrade_url,omitempty"`
	Key        string `json:"key"`
	Limit      int64  `json:"limit"`
}

// WriteError writes a *forfait.LimitError as HTTP 402 in the suite envelope,
// {"error":{"code":"payment_required","message","upgrade_url","key","limit"}},
// and every other error the way httpjson.WriteError does.
func WriteError(w http.ResponseWriter, err error) {
	var refusal *forfait.LimitError
	if !stderrors.As(err, &refusal) {
		httpjson.WriteError(w, err)
		return
	}
	body := limitBody{
		Code:    "payment_required",
		Message: "Your plan's mail account limit is reached. Upgrade to add another one.",
		Key:     refusal.Key,
		Limit:   refusal.Limit,
	}
	if safeUpgradeURL(refusal.UpgradeURL) {
		body.UpgradeURL = refusal.UpgradeURL
	}
	httpjson.WriteJSON(w, http.StatusPaymentRequired, map[string]limitBody{"error": body})
}

// safeUpgradeURL accepts only an absolute https URL with a host and no
// credentials, because the client renders it as a link.
func safeUpgradeURL(raw string) bool {
	parsed, err := url.Parse(raw)
	return err == nil && parsed.Scheme == "https" && parsed.Host != "" && parsed.User == nil
}
