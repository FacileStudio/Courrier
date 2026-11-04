package schemas

import "gorm.io/gorm"

// forfaitStateSchema holds the one document porte/forfait exports after every
// entitlements snapshot, so a restart reads the last known plan values instead
// of failing open again. The check pins the table to a single row.
const forfaitStateSchema = `
CREATE TABLE IF NOT EXISTS forfait_state (
	id         smallint    PRIMARY KEY DEFAULT 1 CHECK (id = 1),
	document   text        NOT NULL,
	updated_at timestamptz NOT NULL DEFAULT now()
);
`

// ensureForfaitState creates the forfait_state table. Idempotent.
func ensureForfaitState(db *gorm.DB) error {
	return db.Exec(forfaitStateSchema).Error
}

// ensureRawTables creates the tables kept as raw DDL rather than GORM models:
// push_devices and forfait_state. Idempotent.
func ensureRawTables(db *gorm.DB) error {
	if err := ensurePushDevices(db); err != nil {
		return err
	}
	return ensureForfaitState(db)
}
