package push

import documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"

var unauthorized = documentation.Error{
	Status: 401, Code: "unauthenticated", Description: "The request carries no valid session.",
}

// Documentation describes the push routes for the API reference.
var Documentation = documentation.Module{
	Name: "push",
	Description: "Push notifications for new inbox mail: the Android app through Firebase Cloud Messaging, " +
		"browsers through Web Push, the desktop app through a live event stream. FCM and Web Push are " +
		"off until their server settings are set.",
	Routes: []documentation.Route{
		{
			Method:  "GET",
			Path:    "/push",
			Summary: "Report the notification channels",
			Description: "Says whether native Android notifications (FCM) are on and returns web_push_key, " +
				"the VAPID public key a browser passes to PushManager.subscribe. The Android app only " +
				"registers for push while native is true.",
			Auth:         "session cookie or bearer token",
			ResponseBody: StatusResponse{},
			Errors:       []documentation.Error{unauthorized},
		},
		{
			Method:  "GET",
			Path:    "/push/stream",
			Summary: "Stream notifications to the desktop app",
			Description: "A text/event-stream the desktop app keeps open: one \"notification\" event whose data is " +
				"the JSON notification, and a comment every 25 seconds. While it is open, the inbox watcher " +
				"polls the caller's inboxes. Notifications live in memory, so this assumes one API replica.",
			Auth:         "session cookie or bearer token",
			ResponseBody: StreamNotification{},
			Errors:       []documentation.Error{unauthorized},
		},
		{
			Method:       "POST",
			Path:         "/push/test",
			Summary:      "Send a test notification",
			Description:  "Sends a generic test notification to every registered device and open desktop app of the caller.",
			Auth:         "session cookie or bearer token",
			RequestBody:  EmptyRequest{},
			ResponseBody: EmptyResponse{},
			Errors:       []documentation.Error{unauthorized},
		},
		{
			Method:  "POST",
			Path:    "/push/devices",
			Summary: "Register a device for notifications",
			Description: "Records an FCM registration token (kind fcm) or a browser Web Push subscription " +
				"(kind web: token is the endpoint, with its p256dh and auth keys). Registering a token " +
				"another account owns moves it to the caller. Answers 204.",
			Auth:         "session cookie or bearer token",
			RequestBody:  DeviceRequest{},
			ResponseBody: EmptyResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "The token or a browser key is missing."},
				unauthorized,
				{Status: 409, Code: "already_exists", Description: "This server cannot send to that kind of device."},
			},
		},
		{
			Method:       "DELETE",
			Path:         "/push/devices",
			Summary:      "Forget a device",
			Description:  "Stops notifying one of the caller's devices. Answers 204.",
			Auth:         "session cookie or bearer token",
			RequestBody:  ForgetDeviceRequest{},
			ResponseBody: EmptyResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "The token is missing."},
				unauthorized,
			},
		},
	},
}
