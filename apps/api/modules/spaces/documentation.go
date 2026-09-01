package spaces

import (
	"net/http"

	documentation "github.com/FacileStudio/Courrier/apps/api/internal/documentation"
)

var spaceID = []documentation.Field{
	{Name: "spaceId", Type: "string", Description: "Space UUID."},
}

var memberPath = []documentation.Field{
	{Name: "spaceId", Type: "string", Description: "Space UUID."},
	{Name: "memberId", Type: "string", Description: "Membership UUID."},
}

var unauthenticated = documentation.Error{Status: 401, Code: "unauthenticated", Description: "No valid session cookie, bearer token or API token."}

var Documentation = documentation.Module{
	Name:        "spaces",
	Description: "Shared workspaces and their membership. Roles are owner, admin and member.",
	Routes: []documentation.Route{
		{
			Method:       "GET",
			Path:         "/spaces",
			Summary:      "List spaces",
			Description:  "Returns only the spaces the caller belongs to, each carrying the caller's own role.",
			Auth:         "bearer token required",
			ResponseBody: SpaceListResponse{},
			Errors:       []documentation.Error{unauthenticated},
		},
		{
			Method:       "POST",
			Path:         "/spaces",
			Summary:      "Create a space",
			Description:  "Creates a space and makes the caller its owner.",
			Auth:         "bearer token required",
			RequestBody:  CreateSpaceRequest{},
			ResponseBody: SpaceResponse{},
			Status:       http.StatusCreated,
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or missing name."},
				unauthenticated,
			},
		},
		{
			Method:       "GET",
			Path:         "/spaces/{spaceId}",
			Summary:      "Get a space",
			Description:  "Returns one space with its members array.",
			Auth:         "bearer token required",
			PathParams:   spaceID,
			ResponseBody: SpaceResponse{},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such space, or the caller is not a member."},
			},
		},
		{
			Method:       "PUT",
			Path:         "/spaces/{spaceId}",
			Summary:      "Update a space",
			Description:  "Renames or re-describes a space. Both fields are optional.",
			Auth:         "bearer token required",
			PathParams:   spaceID,
			RequestBody:  UpdateSpaceRequest{},
			ResponseBody: SpaceResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body."},
				unauthenticated,
				{Status: 403, Code: "permission_denied", Description: "The caller is not an owner or admin of the space."},
				{Status: 404, Code: "not_found", Description: "No such space."},
			},
		},
		{
			Method:       "DELETE",
			Path:         "/spaces/{spaceId}",
			Summary:      "Delete a space",
			Description:  "Deletes the space and every membership in it.",
			Auth:         "bearer token required",
			PathParams:   spaceID,
			ResponseBody: DeleteSpaceResponse{},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 403, Code: "permission_denied", Description: "Only the owner can delete a space."},
				{Status: 404, Code: "not_found", Description: "No such space."},
			},
		},
		{
			Method:       "POST",
			Path:         "/spaces/{spaceId}/leave",
			Summary:      "Leave a space",
			Description:  "Removes the caller's own membership. An owner may leave while another owner remains.",
			Auth:         "bearer token required",
			PathParams:   spaceID,
			ResponseBody: LeaveSpaceResponse{},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "The caller is not a member of that space."},
				{Status: 409, Code: "already_exists", Description: "The caller is the space's only owner."},
			},
		},
		{
			Method:       "GET",
			Path:         "/spaces/{spaceId}/members",
			Summary:      "List space members",
			Description:  "Returns every membership in the space.",
			Auth:         "bearer token required",
			PathParams:   spaceID,
			ResponseBody: MemberListResponse{},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 404, Code: "not_found", Description: "No such space, or the caller is not a member."},
			},
		},
		{
			Method:       "POST",
			Path:         "/spaces/{spaceId}/members",
			Summary:      "Add a member",
			Description:  "Adds an existing user to the space with the given role.",
			Auth:         "bearer token required",
			PathParams:   spaceID,
			RequestBody:  AddMemberRequest{},
			ResponseBody: AddMemberResponse{},
			Status:       http.StatusCreated,
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or unknown role."},
				unauthenticated,
				{Status: 403, Code: "permission_denied", Description: "The caller is not an owner or admin of the space."},
				{Status: 409, Code: "already_exists", Description: "That user is already a member."},
			},
		},
		{
			Method:       "PUT",
			Path:         "/spaces/{spaceId}/members/{memberId}",
			Summary:      "Change a member's role",
			Description:  "Updates one membership's role.",
			Auth:         "bearer token required",
			PathParams:   memberPath,
			RequestBody:  UpdateMemberRequest{},
			ResponseBody: UpdateMemberRoleResponse{},
			Errors: []documentation.Error{
				{Status: 400, Code: "invalid_argument", Description: "Invalid JSON body or unknown role."},
				unauthenticated,
				{Status: 403, Code: "permission_denied", Description: "The caller is not an owner or admin of the space."},
				{Status: 404, Code: "not_found", Description: "No such membership."},
			},
		},
		{
			Method:       "DELETE",
			Path:         "/spaces/{spaceId}/members/{memberId}",
			Summary:      "Remove a member",
			Description:  "Removes someone else's membership from the space.",
			Auth:         "bearer token required",
			PathParams:   memberPath,
			ResponseBody: RemoveMemberResponse{},
			Errors: []documentation.Error{
				unauthenticated,
				{Status: 403, Code: "permission_denied", Description: "The caller is not an owner or admin of the space."},
				{Status: 404, Code: "not_found", Description: "No such membership."},
			},
		},
	},
}
