package settings

// Settings is the settings payload returned to the client.
type Settings struct {
	EncryptionKeySet bool `json:"encryption_key_set"`
}

// UpdateRequest is the body of PUT /settings; every field is optional.
type UpdateRequest struct {
	EncryptionKey *string `json:"encryption_key"`
}

// SettingsResponse wraps the settings payload.
type SettingsResponse struct {
	Settings Settings `json:"settings"`
}
