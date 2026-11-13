package schemas

import "gorm.io/gorm"

// AdoptPorte moves Courrier's credentials onto porte's tables.
//
// issuer is the configured OIDC_ISSUER, and it has to be passed in because it
// is half of the account matching key: porte finds an identity by (provider,
// subject), and the provider is the issuer. Backfilling with the wrong value —
// or with a placeholder — would leave every existing SSO user unmatched, which
// degrades silently to the email fallback rather than failing.
//
// An empty issuer skips only the identity and picture backfills. The sessions,
// the API tokens and the passwords still move, because they are what keeps
// people signed in and able to sign in again.
func AdoptPorte(db *gorm.DB, issuer string) error {
	statements := []string{rekeyLocalIdentities, carryLegacySessionsOver, carryAPITokensOver, adoptExistingPasswords}
	for _, statement := range statements {
		if err := db.Exec(statement).Error; err != nil {
			return err
		}
	}
	if issuer == "" {
		return nil
	}
	if err := db.Exec(adoptOIDCSubjects, issuer).Error; err != nil {
		return err
	}
	return db.Exec(backfillIssuerPictures, issuer, issuer).Error
}

// rekeyLocalIdentities is porte v0.3.0 re-keying password identities off the
// email address and onto the account id. porte's tables come from the
// baseline migration, a private copy of porte/pg's Schema with porte_users
// left out and the foreign keys pointed at Courrier's users, so this UPDATE is
// copied rather than received: without it every existing password login
// would answer "invalid credentials" to the right password.
//
// It is idempotent, since after it runs the predicate is false for every row.
const rekeyLocalIdentities = `
UPDATE porte_identities SET subject = user_id::text
 WHERE provider = 'local' AND subject <> user_id::text;
`

// carryLegacySessionsOver moves the pre-porte sessions table across instead of
// dropping it, so adopting porte does not sign every existing user out.
//
// Both tables store the SHA-256 hex of a 32-byte token and never the token
// itself, so a row copied here keeps authenticating the credential already in
// somebody's browser — and porte's AcceptLegacyCookie keeps reading it from
// the `session` cookie this app has always set.
//
// last_used_at is stamped now rather than copied from created_at: porte
// retires a browser session idle for seven days and the old table recorded no
// use at all, so carrying created_at over would sign out everyone who last
// signed in more than a week ago, on the deploy meant to keep them.
const carryLegacySessionsOver = `
DO $$
BEGIN
	IF to_regclass('sessions') IS NOT NULL THEN
		INSERT INTO porte_sessions (token_hash, user_id, created_at, last_used_at, expires_at)
		SELECT s.token, s.user_id, s.created_at, now(), s.expires_at
		  FROM sessions s JOIN users u ON u.id = s.user_id
		 WHERE s.expires_at > now()
		ON CONFLICT (token_hash) DO NOTHING;
		DROP TABLE sessions;
	END IF;
END
$$;
`

// carryAPITokensOver folds the separate api_tokens table into porte's sessions
// as labelled rows.
//
// porte.Session.Label exists for precisely this: what Courrier called an API
// token is a session with a name on it and no expiry, which is one mechanism
// more than the problem needs and a second branch in the authentication path.
// A null expires_at never expires, and porte's idle window deliberately does
// not apply to a labelled session, so a token wired into a script keeps
// working.
const carryAPITokensOver = `
DO $$
BEGIN
	IF to_regclass('api_tokens') IS NOT NULL THEN
		INSERT INTO porte_sessions (token_hash, user_id, label, created_at, last_used_at, expires_at)
		SELECT t.token, t.user_id, COALESCE(NULLIF(t.name, ''), 'API token'), t.created_at, now(), NULL
		  FROM api_tokens t JOIN users u ON u.id = t.user_id
		ON CONFLICT (token_hash) DO NOTHING;
		DROP TABLE api_tokens;
	END IF;
END
$$;
`

// adoptExistingPasswords moves the argon2 hashes from users.password_hash into
// the identity rows porte/local reads.
//
// Without it the deploy silently ends password login for every existing
// account: the hash is still in the users table, nothing reads it there any
// more, and the login form answers "invalid credentials" to a correct
// password with no error anywhere. The hashes are byte-identical — porte/local
// uses the parameters this app already used — so the move is a copy and nobody
// resets anything.
//
// The subject is the account id, which is what porte.LocalSubject returns.
// v0.2 keyed it on the lowercased address and that key moved every time
// somebody edited their profile, so this app carried an UPDATE against
// porte_identities to chase it. The id does not move, so the chase is gone.
//
// It has to agree with the re-key above, which runs first: an INSERT still
// writing the address would re-create an email-keyed row beside the id-keyed
// one on every boot, and the next boot's re-key would collide with the
// primary key and fail the migration for good.
//
// users.password_hash is deliberately left in place. Blanking it in the same
// deploy makes the change unrollbackable for the sake of tidiness, and a
// column nothing reads can be dropped on any later day.
const adoptExistingPasswords = `
INSERT INTO porte_identities (user_id, provider, subject, password_hash)
SELECT id, 'local', id::text, password_hash
  FROM users
 WHERE coalesce(password_hash, '') <> ''
ON CONFLICT (provider, subject) DO NOTHING;
`

// adoptOIDCSubjects moves the federated identity off the user row and into the
// identity table porte reads.
//
// Without it nothing breaks loudly: porte finds no identity for (issuer,
// subject), falls back to matching the verified email, finds the same user and
// links it on the next login. But the email fallback is the weaker path by
// design — it is the one an identity provider that lets a user set any address
// can abuse — and relying on it for every existing account, on the one deploy
// where it would go unnoticed, is not a trade worth making.
//
// The OIDC tokens deliberately do **not** come across. Courrier stores them
// encrypted with ENCRYPTION_KEY and porte stores them as it will send them, so
// copying the ciphertext across would hand porte a refresh token that is not
// one: the first profile sync would fail, clear both, and look like the
// provider had revoked them. They are re-issued in full on the user's next SSO
// login, and the only thing missing until then is a background profile
// refresh.
const adoptOIDCSubjects = `
INSERT INTO porte_identities (user_id, provider, subject)
SELECT id, ?, oidc_subject
  FROM users
 WHERE oidc_subject IS NOT NULL AND oidc_subject <> ''
ON CONFLICT (provider, subject) DO NOTHING;
`

// backfillIssuerPictures gives every federated user without a stored picture
// the identity provider's stable avatar URL, <issuer>/avatars/<subject>.
//
// Registre serves that URL for every human and answers 404 when there is no
// picture, which the client draws as initials. Before this, a user who signed
// in before Registre emitted the claim showed initials until their next
// sign-in. It only fills empty values, so it is idempotent and never
// overwrites what the last sign-in or profile sync wrote.
const backfillIssuerPictures = `
UPDATE users u
   SET oidc_picture_url = rtrim(?, '/') || '/avatars/' || i.subject
  FROM porte_identities i
 WHERE i.user_id = u.id
   AND i.provider = ?
   AND coalesce(u.oidc_picture_url, '') = '';
`
