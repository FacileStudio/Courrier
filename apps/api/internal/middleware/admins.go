package middleware

import "strings"

// Admins is the COURRIER_ADMINS allow-list: the emails that may read every
// account on the instance. Courrier has no admin role of its own and any
// Registre account can sign in, so an empty list grants that to nobody.
type Admins []string

// Allows reports whether email is on the list, ignoring case and surrounding space.
func (admins Admins) Allows(email string) bool {
	want := strings.ToLower(strings.TrimSpace(email))
	if want == "" {
		return false
	}
	for _, admin := range admins {
		if strings.ToLower(strings.TrimSpace(admin)) == want {
			return true
		}
	}
	return false
}
