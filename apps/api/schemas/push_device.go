package schemas

import "gorm.io/gorm"

// pushDevicesSchema is the table internal/push stores its devices in: one row
// per FCM registration token or Web Push endpoint. The token is the key, so a
// phone or browser shared between two accounts belongs to whoever registered
// it last.
const pushDevicesSchema = `
CREATE TABLE IF NOT EXISTS push_devices (
	token      text        PRIMARY KEY,
	user_id    bigint      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	kind       text        NOT NULL,
	p256dh     text        NOT NULL DEFAULT '',
	auth       text        NOT NULL DEFAULT '',
	created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS push_devices_user_idx ON push_devices (user_id);
`

// ensurePushDevices creates the push_devices table. Idempotent.
func ensurePushDevices(db *gorm.DB) error {
	return db.Exec(pushDevicesSchema).Error
}
