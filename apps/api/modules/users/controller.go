package users

import (
	"bytes"
	"context"
	stderrors "errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FacileStudio/Courrier/apps/api/internal/authcontext"
	"github.com/FacileStudio/porte"
	"github.com/FacileStudio/tronc/errors"
)

// Controller is the users module's HTTP adapter.
type Controller struct {
	service *Service
}

func newController(service *Service) *Controller {
	return &Controller{service: service}
}

func (controller *Controller) list(context context.Context) (*ListResponse, error) {
	if _, ok := authcontext.IdentityFromContext(context); !ok {
		return nil, errors.Unauthorized("missing auth")
	}

	users, err := controller.service.listUsers(context)
	if err != nil {
		return nil, err
	}

	return &ListResponse{Users: users}, nil
}

func (controller *Controller) get(context context.Context, userID string) (*MeResponse, error) {
	user, err := controller.service.getUser(context, userID)
	if err != nil {
		return nil, err
	}
	return &MeResponse{User: *user}, nil
}

func (controller *Controller) me(context context.Context) (*MeResponse, error) {
	identity, ok := authcontext.IdentityFromContext(context)
	if !ok {
		return nil, errors.Unauthorized("missing auth")
	}

	user, err := controller.service.getUser(context, identity.UserID)
	if err != nil {
		return nil, err
	}

	if user.Email == "" {
		user.Email = identity.Email
	}

	return &MeResponse{User: *user}, nil
}

func (controller *Controller) updateMe(context context.Context, w http.ResponseWriter, request *http.Request, req *UpdateRequest) (*MeResponse, error) {
	identity, ok := authcontext.IdentityFromContext(context)
	if !ok {
		return nil, errors.Unauthorized("missing auth")
	}

	var name *string
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if len(trimmed) > 80 {
			return nil, errors.Invalid("name must be at most 80 characters")
		}
		name = &trimmed
	}

	var email *string
	if req.Email != nil {
		normalized := strings.TrimSpace(strings.ToLower(*req.Email))
		if normalized == "" || !strings.Contains(normalized, "@") {
			return nil, errors.Invalid("invalid email")
		}
		email = &normalized
	}

	var password *string
	if req.Password != nil {
		if len(*req.Password) < 8 {
			return nil, errors.Invalid("password must be at least 8 characters")
		}
		password = req.Password
	}

	if name == nil && email == nil && password == nil {
		return nil, errors.Invalid("at least one field must be provided")
	}

	if password != nil {
		if err := controller.changePassword(context, w, request, identity.UserID, req, *password); err != nil {
			return nil, err
		}
	}

	user, err := controller.service.updateUser(context, identity.UserID, name, email)
	if err != nil {
		return nil, err
	}

	return &MeResponse{User: *user}, nil
}

// changePassword picks between porte's two password writes.
//
// They are two calls and not one because only one of them is safe to make
// with nothing but a session: SetPassword gives a first password to an account
// that has none, and porte refuses it with ErrPasswordSet once there is one.
// Replacing a password goes through ChangePassword, which confirms the current
// one — OWASP ASVS puts that at L1 (v4 2.1.6, v5 6.2.3) — then ends the
// account's other logins and rotates this caller's session through w.
//
// A blank current password counts as none given, so the answer is the one that
// says what is missing rather than "invalid credentials".
func (controller *Controller) changePassword(context context.Context, w http.ResponseWriter, request *http.Request, userID string, req *UpdateRequest, password string) error {
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return errors.Internal("failed to parse user id", err)
	}

	current := ""
	if req.CurrentPassword != nil {
		current = strings.TrimSpace(*req.CurrentPassword)
	}
	if current != "" {
		_, _, err := controller.service.tokens.ChangePassword(context, w, request, id, current, password)
		return err
	}

	err = controller.service.tokens.SetPassword(context, id, password)
	if stderrors.Is(err, porte.ErrPasswordSet) {
		return errors.Invalid("current password is required to change your password")
	}
	return err
}

func (controller *Controller) deleteAvatar(context context.Context) (*MeResponse, error) {
	identity, ok := authcontext.IdentityFromContext(context)
	if !ok {
		return nil, errors.Unauthorized("missing auth")
	}
	user, err := controller.service.clearAvatar(context, identity.UserID)
	if err != nil {
		return nil, err
	}
	return &MeResponse{User: *user}, nil
}

func (controller *Controller) uploadAvatar(context context.Context, request *http.Request) (*MeResponse, error) {
	identity, ok := authcontext.IdentityFromContext(context)
	if !ok {
		return nil, errors.Unauthorized("missing auth")
	}

	if err := request.ParseMultipartForm(5 << 20); err != nil {
		return nil, errors.TooLarge("avatar file is too large")
	}

	file, _, err := request.FormFile("avatar")
	if err != nil {
		return nil, errors.Invalid("avatar file is required")
	}
	defer file.Close()

	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return nil, errors.Internal("failed to read avatar file", err)
	}

	contentType := http.DetectContentType(header[:n])
	user, err := controller.service.storeAvatar(context, identity.UserID, io.MultiReader(bytes.NewReader(header[:n]), file), contentType)
	if err != nil {
		return nil, err
	}

	return &MeResponse{User: *user}, nil
}

func (controller *Controller) getApiToken(context context.Context) (*ApiTokenStatusResponse, error) {
	identity, ok := authcontext.IdentityFromContext(context)
	if !ok {
		return nil, errors.Unauthorized("missing auth")
	}
	record, err := controller.service.getApiToken(context, identity.UserID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return &ApiTokenStatusResponse{HasToken: false}, nil
	}
	return &ApiTokenStatusResponse{
		HasToken:  true,
		Name:      record.Label,
		CreatedAt: record.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

func (controller *Controller) createApiToken(context context.Context, req *CreateApiTokenRequest) (*ApiTokenResponse, error) {
	identity, ok := authcontext.IdentityFromContext(context)
	if !ok {
		return nil, errors.Unauthorized("missing auth")
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = "CLI"
	}
	rawToken, record, err := controller.service.createApiToken(context, identity.UserID, name)
	if err != nil {
		return nil, err
	}
	return &ApiTokenResponse{
		Token:     rawToken,
		Name:      record.Label,
		CreatedAt: record.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

func (controller *Controller) deleteApiToken(context context.Context) error {
	identity, ok := authcontext.IdentityFromContext(context)
	if !ok {
		return errors.Unauthorized("missing auth")
	}
	return controller.service.deleteApiToken(context, identity.UserID)
}
