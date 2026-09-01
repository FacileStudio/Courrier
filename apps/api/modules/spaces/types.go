package spaces

// CreateSpaceRequest is the body of POST /spaces.
type CreateSpaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateSpaceRequest is the body of PUT /spaces/{id}; every field is optional.
type UpdateSpaceRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// AddMemberRequest is the body of POST /spaces/{id}/members.
type AddMemberRequest struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
}

// UpdateMemberRequest is the body of PUT /spaces/{id}/members/{memberId}.
type UpdateMemberRequest struct {
	Role string `json:"role"`
}

// SpaceResponse is the space shape returned to the client, including the
// caller's role and, for detail views, its members.
type SpaceResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Role        string           `json:"role"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
	Members     []MemberResponse `json:"members,omitempty"`
}

// MemberResponse is a space membership as returned to the client.
type MemberResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
}

type SpaceListResponse struct {
	Spaces []SpaceResponse `json:"spaces"`
}

type DeleteSpaceResponse struct {
	Deleted bool `json:"deleted"`
}

type LeaveSpaceResponse struct {
	Left bool `json:"left"`
}

type MemberListResponse struct {
	Members []MemberResponse `json:"members"`
}

type AddMemberResponse struct {
	ID      string `json:"id"`
	SpaceID string `json:"space_id"`
	UserID  int64  `json:"user_id"`
	Role    string `json:"role"`
}

type UpdateMemberRoleResponse struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type RemoveMemberResponse struct {
	Removed bool `json:"removed"`
}
