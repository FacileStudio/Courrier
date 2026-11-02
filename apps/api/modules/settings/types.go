package settings

// Settings is the settings payload returned to the client.
type Settings struct {
	EncryptionKeySet bool `json:"encryption_key_set"`
}

// SettingsResponse wraps the settings payload.
type SettingsResponse struct {
	Settings Settings `json:"settings"`
}
