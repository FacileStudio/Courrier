package main

import (
	"github.com/FacileStudio/Courrier/apps/api/internal/bulle"
	"github.com/FacileStudio/Courrier/apps/api/internal/env"
	"github.com/FacileStudio/Courrier/apps/api/modules/assist"
)

// newAssist builds the Bulle writing helper. It stays off (every call answers
// "off") until BULLE_URL and the Registre client are both configured.
func newAssist(appEnv env.Config, ids bulle.Identities) *assist.Service {
	if appEnv.OIDC == nil {
		return assist.NewService(nil, "")
	}
	return assist.NewService(bulle.New(bulle.Config{
		BulleURL:     appEnv.BulleURL,
		Issuer:       appEnv.OIDC.Issuer,
		ClientID:     appEnv.OIDC.ClientID,
		ClientSecret: appEnv.OIDC.ClientSecret,
	}, ids, nil), appEnv.BulleURL)
}
