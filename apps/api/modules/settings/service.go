package settings

// Service reports instance settings that come from the server environment.
type Service struct {
	encryptionKeySet bool
	controller       *Controller
}

// NewService builds a settings Service. encryptionKeySet says whether
// ENCRYPTION_KEY was loaded at startup.
func NewService(encryptionKeySet bool) *Service {
	service := &Service{encryptionKeySet: encryptionKeySet}
	service.controller = newController(service)
	return service
}

func (service *Service) getSettings() *Settings {
	return &Settings{EncryptionKeySet: service.encryptionKeySet}
}
