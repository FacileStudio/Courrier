package push

import (
	"context"

	"gorm.io/gorm"
)

// PGDevices stores devices in the push_devices table.
type PGDevices struct{ db *gorm.DB }

// NewPGDevices returns a DeviceStore backed by Postgres.
func NewPGDevices(db *gorm.DB) *PGDevices { return &PGDevices{db: db} }

// SaveDevice records a device, moving it to this user when another account
// registered the same token on a shared phone or browser.
func (p *PGDevices) SaveDevice(ctx context.Context, d Device) error {
	return p.db.WithContext(ctx).Exec(
		`INSERT INTO push_devices (token, user_id, kind, p256dh, auth) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT (token) DO UPDATE SET user_id = EXCLUDED.user_id, kind = EXCLUDED.kind,
		p256dh = EXCLUDED.p256dh, auth = EXCLUDED.auth`,
		d.Token, d.UserID, d.Kind, d.P256dh, d.Auth,
	).Error
}

// DeleteDevice removes one of the user's devices.
func (p *PGDevices) DeleteDevice(ctx context.Context, userID int64, token string) error {
	return p.db.WithContext(ctx).Exec(
		"DELETE FROM push_devices WHERE user_id = ? AND token = ?", userID, token,
	).Error
}

// Devices lists the user's devices.
func (p *PGDevices) Devices(ctx context.Context, userID int64) ([]Device, error) {
	var out []Device
	err := p.db.WithContext(ctx).Raw(
		"SELECT user_id, kind, token, p256dh, auth FROM push_devices WHERE user_id = ?", userID,
	).Scan(&out).Error
	return out, err
}
