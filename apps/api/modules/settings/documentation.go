package settings

import documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"

var Documentation = documentation.Module{
	Name:        "settings",
	Description: "Application settings management.",
	Routes: []documentation.Route{
		{
			Method:       "GET",
			Path:         "/settings",
			Summary:      "Return instance settings",
			Description:  "Reports whether ENCRYPTION_KEY was loaded at startup.",
			Auth:         "bearer token required",
			ResponseBody: SettingsResponse{},
			Errors: []documentation.Error{
				{Status: 401, Code: "unauthenticated", Description: "Authorization header is missing or invalid."},
			},
		},
	},
}
