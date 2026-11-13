-- +goose Up
-- Rows written outside GORM got NULL where GORM used to fill the column itself. The
-- database now fills the timestamps GORM stamped at insert, and the users columns
-- kit.Store does not name get the values modules/auth's fillNewAccount writes.

ALTER TABLE accounts ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE accounts ALTER COLUMN updated_at SET DEFAULT now();
ALTER TABLE email_templates ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE email_templates ALTER COLUMN updated_at SET DEFAULT now();
ALTER TABLE emails ALTER COLUMN cached_at SET DEFAULT now();
ALTER TABLE sender_rules ALTER COLUMN updated_at SET DEFAULT now();
ALTER TABLE space_members ALTER COLUMN joined_at SET DEFAULT now();
ALTER TABLE spaces ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE spaces ALTER COLUMN updated_at SET DEFAULT now();

ALTER TABLE users
    ALTER COLUMN created_at SET DEFAULT now(),
    ALTER COLUMN avatar_url SET DEFAULT '',
    ALTER COLUMN avatar_source SET DEFAULT '',
    ALTER COLUMN avatar_upload_path SET DEFAULT '',
    ALTER COLUMN password_hash SET DEFAULT '',
    ALTER COLUMN oidc_access_token SET DEFAULT '',
    ALTER COLUMN oidc_refresh_token SET DEFAULT '',
    ALTER COLUMN oidc_token_expiry SET DEFAULT '0001-01-01 00:00:00+00',
    ALTER COLUMN profile_synced_at SET DEFAULT '0001-01-01 00:00:00+00';

-- +goose Down
ALTER TABLE users
    ALTER COLUMN created_at DROP DEFAULT,
    ALTER COLUMN avatar_url DROP DEFAULT,
    ALTER COLUMN avatar_source DROP DEFAULT,
    ALTER COLUMN avatar_upload_path DROP DEFAULT,
    ALTER COLUMN password_hash DROP DEFAULT,
    ALTER COLUMN oidc_access_token DROP DEFAULT,
    ALTER COLUMN oidc_refresh_token DROP DEFAULT,
    ALTER COLUMN oidc_token_expiry DROP DEFAULT,
    ALTER COLUMN profile_synced_at DROP DEFAULT;

ALTER TABLE spaces ALTER COLUMN updated_at DROP DEFAULT;
ALTER TABLE spaces ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE space_members ALTER COLUMN joined_at DROP DEFAULT;
ALTER TABLE sender_rules ALTER COLUMN updated_at DROP DEFAULT;
ALTER TABLE emails ALTER COLUMN cached_at DROP DEFAULT;
ALTER TABLE email_templates ALTER COLUMN updated_at DROP DEFAULT;
ALTER TABLE email_templates ALTER COLUMN created_at DROP DEFAULT;
ALTER TABLE accounts ALTER COLUMN updated_at DROP DEFAULT;
ALTER TABLE accounts ALTER COLUMN created_at DROP DEFAULT;
