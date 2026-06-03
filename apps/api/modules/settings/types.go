package settings

type Settings struct {
	NookPoolURL     string `json:"nook_pool_url"`
	NookPoolSecret  string `json:"nook_pool_secret"`
	NookPoolEnabled bool   `json:"nook_pool_enabled"`
}

type UpdateRequest struct {
	NookPoolURL     string `json:"nook_pool_url"`
	NookPoolSecret  string `json:"nook_pool_secret"`
	NookPoolEnabled bool   `json:"nook_pool_enabled"`
}

type SettingsResponse struct {
	Settings Settings `json:"settings"`
}
