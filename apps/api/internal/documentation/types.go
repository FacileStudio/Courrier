package docs

import "github.com/FacileStudio/tronc/apiref"

type (
	// Response is the top-level registry the API reference generator renders.
	Response = apiref.Registry
	// Module is a documented API module in the reference.
	Module = apiref.Module
	// Route is a documented endpoint in the reference.
	Route = apiref.Route
	// Field is a documented request or response field.
	Field = apiref.Field
	// Error is a documented error response shape.
	Error = apiref.Error
)
