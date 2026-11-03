// Package push delivers generic notifications to a user's phones and browsers:
// native Android through Firebase Cloud Messaging and browsers through Web Push.
package push

import (
	"context"
	"errors"
	"time"
)

const (
	sendTimeout = 5 * time.Second

	KindFCM = "fcm"
	KindWeb = "web"
)

// ErrUnsupported means the server cannot deliver to that kind of device.
var ErrUnsupported = errors.New("this server cannot send to that kind of device")

// Notification is a generic message. Title and Body must never carry mail
// content: a sender name is the most a notification may say.
type Notification struct {
	Title string
	Body  string
	Path  string
}

// Device is one place a user receives notifications: an FCM registration
// token, or a Web Push endpoint with its encryption keys.
type Device struct {
	UserID int64
	Kind   string
	Token  string
	P256dh string
	Auth   string
}

// DeviceStore keeps the devices of each user.
type DeviceStore interface {
	SaveDevice(ctx context.Context, d Device) error
	DeleteDevice(ctx context.Context, userID int64, token string) error
	Devices(ctx context.Context, userID int64) ([]Device, error)
}
