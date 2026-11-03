package push

import (
	"context"
	"errors"
	"log/slog"
)

// HubConfig lists the delivery channels. FCM and Web may each be nil.
type HubConfig struct {
	FCM     *FCM
	Web     *WebPush
	Devices DeviceStore
	Logger  *slog.Logger
}

// Hub fans one notification out to every device of the user.
// A nil Hub is a safe no-op.
type Hub struct{ cfg HubConfig }

// NewHub returns a Hub over the configured channels.
func NewHub(cfg HubConfig) *Hub {
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	return &Hub{cfg: cfg}
}

// Native reports whether native Android notifications are on.
func (h *Hub) Native() bool { return h != nil && h.cfg.FCM != nil && h.cfg.Devices != nil }

// WebPushKey returns the VAPID public key browsers subscribe with, or "".
func (h *Hub) WebPushKey() string {
	if h == nil || h.cfg.Web == nil || h.cfg.Devices == nil {
		return ""
	}
	return h.cfg.Web.publicKey
}

// Enabled reports whether any channel delivers notifications.
func (h *Hub) Enabled() bool { return h.Native() || h.WebPushKey() != "" }

// Register records a device when the server can send to its kind.
func (h *Hub) Register(ctx context.Context, d Device) error {
	switch {
	case d.Kind == KindFCM && h.Native():
	case d.Kind == KindWeb && h.WebPushKey() != "":
	default:
		return ErrUnsupported
	}
	return h.cfg.Devices.SaveDevice(ctx, d)
}

// Unregister forgets one of the user's devices.
func (h *Hub) Unregister(ctx context.Context, userID int64, token string) error {
	if h == nil || h.cfg.Devices == nil {
		return nil
	}
	return h.cfg.Devices.DeleteDevice(ctx, userID, token)
}

// Notify delivers to every device of the user. Failures are logged and never
// returned; a device the push service reports as gone is forgotten.
func (h *Hub) Notify(ctx context.Context, userID int64, n Notification) {
	if !h.Enabled() {
		return
	}
	devices, err := h.cfg.Devices.Devices(ctx, userID)
	if err != nil {
		h.cfg.Logger.Warn("listing push devices failed", slog.Any("error", err))
		return
	}
	for _, d := range devices {
		h.deliver(ctx, d, n)
	}
}

func (h *Hub) deliver(ctx context.Context, d Device, n Notification) {
	var err error
	switch {
	case d.Kind == KindFCM && h.cfg.FCM != nil:
		err = h.cfg.FCM.send(ctx, d.Token, n)
	case d.Kind == KindWeb && h.cfg.Web != nil:
		err = h.cfg.Web.send(ctx, d, n)
	default:
		return
	}
	if errors.Is(err, ErrGone) {
		_ = h.cfg.Devices.DeleteDevice(ctx, d.UserID, d.Token)
		return
	}
	if err != nil {
		h.cfg.Logger.Warn("push delivery failed", slog.String("kind", d.Kind), slog.Any("error", err))
	}
}
