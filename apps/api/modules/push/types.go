package push

// StatusResponse says which channels deliver notifications to this server's
// users: native is FCM for the Android app, web_push_key is the VAPID public
// key browsers subscribe with and is empty while Web Push is off.
type StatusResponse struct {
	Enabled    bool   `json:"enabled"`
	Native     bool   `json:"native"`
	WebPushKey string `json:"web_push_key,omitempty"`
}

// DeviceRequest registers a device: an FCM token (kind "fcm"), or a browser
// subscription endpoint with its keys (kind "web").
type DeviceRequest struct {
	Kind   string `json:"kind"`
	Token  string `json:"token"`
	P256dh string `json:"p256dh,omitempty"`
	Auth   string `json:"auth,omitempty"`
}

// ForgetDeviceRequest names the device to stop notifying.
type ForgetDeviceRequest struct {
	Token string `json:"token"`
}

// EmptyRequest documents a route that takes no body.
type EmptyRequest struct{}

// EmptyResponse documents a 204 answer.
type EmptyResponse struct{}
