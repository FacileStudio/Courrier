package settings

// Controller is the settings module's HTTP adapter.
type Controller struct {
	service *Service
}

func newController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) getSettings() *SettingsResponse {
	return &SettingsResponse{Settings: *c.service.getSettings()}
}
