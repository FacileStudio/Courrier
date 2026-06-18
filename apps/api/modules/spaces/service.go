package spaces

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"api/internal/errors"
	"api/schemas"

	"gorm.io/gorm"
)

type Service struct {
	orm *gorm.DB
}

func NewService(orm *gorm.DB) *Service {
	return &Service{orm: orm}
}

func (s *Service) Create(ctx context.Context, userID int64, req CreateSpaceRequest) (schemas.Space, string, error) {
	if req.Name == "" {
		return schemas.Space{}, "", errors.Invalid("name is required")
	}

	space := schemas.Space{
		Name:        req.Name,
		Description: req.Description,
	}

	err := s.orm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&space).Error; err != nil {
			return fmt.Errorf("create space: %w", err)
		}

		member := schemas.SpaceMember{
			SpaceID: space.ID,
			UserID:  userID,
			Role:    schemas.SpaceRoleOwner,
		}
		if err := tx.Create(&member).Error; err != nil {
			return fmt.Errorf("create owner member: %w", err)
		}

		return nil
	})
	if err != nil {
		return schemas.Space{}, "", errors.Internal("failed to create space", err)
	}

	return space, schemas.SpaceRoleOwner, nil
}

func (s *Service) List(ctx context.Context, userID int64) ([]schemas.Space, map[string]string, error) {
	var members []schemas.SpaceMember
	if err := s.orm.WithContext(ctx).Where("user_id = ?", userID).Find(&members).Error; err != nil {
		return nil, nil, errors.Internal("failed to list memberships", err)
	}

	if len(members) == 0 {
		return []schemas.Space{}, map[string]string{}, nil
	}

	spaceIDs := make([]string, len(members))
	roleMap := make(map[string]string, len(members))
	for i, m := range members {
		spaceIDs[i] = m.SpaceID
		roleMap[m.SpaceID] = m.Role
	}

	var spaceList []schemas.Space
	if err := s.orm.WithContext(ctx).Where("id IN ?", spaceIDs).Order("name ASC").Find(&spaceList).Error; err != nil {
		return nil, nil, errors.Internal("failed to list spaces", err)
	}

	return spaceList, roleMap, nil
}

func (s *Service) Get(ctx context.Context, userID int64, spaceID string) (schemas.Space, string, error) {
	role, err := s.memberRole(ctx, spaceID, userID)
	if err != nil {
		return schemas.Space{}, "", err
	}

	var space schemas.Space
	if err := s.orm.WithContext(ctx).Preload("Members").Where("id = ?", spaceID).First(&space).Error; err != nil {
		return schemas.Space{}, "", errors.NotFound("space not found")
	}

	return space, role, nil
}

func (s *Service) Update(ctx context.Context, userID int64, spaceID string, req UpdateSpaceRequest) (schemas.Space, error) {
	role, err := s.memberRole(ctx, spaceID, userID)
	if err != nil {
		return schemas.Space{}, err
	}
	if role != schemas.SpaceRoleOwner && role != schemas.SpaceRoleAdmin {
		return schemas.Space{}, errors.Forbidden("only owners and admins can update a space")
	}

	var space schemas.Space
	if err := s.orm.WithContext(ctx).Where("id = ?", spaceID).First(&space).Error; err != nil {
		return schemas.Space{}, errors.NotFound("space not found")
	}

	if req.Name != nil {
		space.Name = *req.Name
	}
	if req.Description != nil {
		space.Description = *req.Description
	}

	if err := s.orm.WithContext(ctx).Save(&space).Error; err != nil {
		return schemas.Space{}, errors.Internal("failed to update space", err)
	}

	return space, nil
}

func (s *Service) Delete(ctx context.Context, userID int64, spaceID string) error {
	role, err := s.memberRole(ctx, spaceID, userID)
	if err != nil {
		return err
	}
	if role != schemas.SpaceRoleOwner {
		return errors.Forbidden("only the owner can delete a space")
	}

	result := s.orm.WithContext(ctx).Where("id = ?", spaceID).Delete(&schemas.Space{})
	if result.Error != nil {
		return errors.Internal("failed to delete space", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("space not found")
	}

	return nil
}

func (s *Service) AddMember(ctx context.Context, actorID int64, spaceID string, req AddMemberRequest) (schemas.SpaceMember, error) {
	role, err := s.memberRole(ctx, spaceID, actorID)
	if err != nil {
		return schemas.SpaceMember{}, err
	}
	if role != schemas.SpaceRoleOwner && role != schemas.SpaceRoleAdmin {
		return schemas.SpaceMember{}, errors.Forbidden("only owners and admins can add members")
	}

	memberRole := req.Role
	if memberRole == "" {
		memberRole = schemas.SpaceRoleMember
	}
	if memberRole != schemas.SpaceRoleAdmin && memberRole != schemas.SpaceRoleMember {
		return schemas.SpaceMember{}, errors.Invalid("role must be admin or member")
	}

	var existing schemas.SpaceMember
	if err := s.orm.WithContext(ctx).Where("space_id = ? AND user_id = ?", spaceID, req.UserID).First(&existing).Error; err == nil {
		return schemas.SpaceMember{}, errors.Conflict("user is already a member")
	}

	var user schemas.User
	if err := s.orm.WithContext(ctx).Where("id = ?", req.UserID).First(&user).Error; err != nil {
		return schemas.SpaceMember{}, errors.NotFound("user not found")
	}

	member := schemas.SpaceMember{
		SpaceID: spaceID,
		UserID:  req.UserID,
		Role:    memberRole,
	}
	if err := s.orm.WithContext(ctx).Create(&member).Error; err != nil {
		return schemas.SpaceMember{}, errors.Internal("failed to add member", err)
	}

	return member, nil
}

func (s *Service) UpdateMember(ctx context.Context, actorID int64, spaceID string, memberID string, req UpdateMemberRequest) (schemas.SpaceMember, error) {
	role, err := s.memberRole(ctx, spaceID, actorID)
	if err != nil {
		return schemas.SpaceMember{}, err
	}
	if role != schemas.SpaceRoleOwner && role != schemas.SpaceRoleAdmin {
		return schemas.SpaceMember{}, errors.Forbidden("only owners and admins can update members")
	}

	if req.Role != schemas.SpaceRoleAdmin && req.Role != schemas.SpaceRoleMember {
		return schemas.SpaceMember{}, errors.Invalid("role must be admin or member")
	}

	var member schemas.SpaceMember
	if err := s.orm.WithContext(ctx).Where("id = ? AND space_id = ?", memberID, spaceID).First(&member).Error; err != nil {
		return schemas.SpaceMember{}, errors.NotFound("member not found")
	}

	if member.Role == schemas.SpaceRoleOwner {
		return schemas.SpaceMember{}, errors.Forbidden("cannot change the owner's role")
	}

	member.Role = req.Role
	if err := s.orm.WithContext(ctx).Save(&member).Error; err != nil {
		return schemas.SpaceMember{}, errors.Internal("failed to update member", err)
	}

	return member, nil
}

func (s *Service) RemoveMember(ctx context.Context, actorID int64, spaceID string, memberID string) error {
	role, err := s.memberRole(ctx, spaceID, actorID)
	if err != nil {
		return err
	}
	if role != schemas.SpaceRoleOwner && role != schemas.SpaceRoleAdmin {
		return errors.Forbidden("only owners and admins can remove members")
	}

	var member schemas.SpaceMember
	if err := s.orm.WithContext(ctx).Where("id = ? AND space_id = ?", memberID, spaceID).First(&member).Error; err != nil {
		return errors.NotFound("member not found")
	}

	if member.Role == schemas.SpaceRoleOwner {
		return errors.Forbidden("cannot remove the owner")
	}

	if err := s.orm.WithContext(ctx).Delete(&member).Error; err != nil {
		return errors.Internal("failed to remove member", err)
	}

	return nil
}

func (s *Service) Leave(ctx context.Context, userID int64, spaceID string) error {
	var member schemas.SpaceMember
	if err := s.orm.WithContext(ctx).Where("space_id = ? AND user_id = ?", spaceID, userID).First(&member).Error; err != nil {
		return errors.NotFound("not a member of this space")
	}

	if member.Role == schemas.SpaceRoleOwner {
		return errors.Forbidden("the owner cannot leave; delete the space or transfer ownership")
	}

	if err := s.orm.WithContext(ctx).Delete(&member).Error; err != nil {
		return errors.Internal("failed to leave space", err)
	}

	return nil
}

func (s *Service) ListMembers(ctx context.Context, userID int64, spaceID string) ([]MemberResponse, error) {
	if _, err := s.memberRole(ctx, spaceID, userID); err != nil {
		return nil, err
	}

	var members []schemas.SpaceMember
	if err := s.orm.WithContext(ctx).Where("space_id = ?", spaceID).Find(&members).Error; err != nil {
		return nil, errors.Internal("failed to list members", err)
	}

	userIDs := make([]int64, len(members))
	for i, m := range members {
		userIDs[i] = m.UserID
	}

	var users []schemas.User
	if err := s.orm.WithContext(ctx).Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, errors.Internal("failed to load users", err)
	}

	userMap := make(map[int64]schemas.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	resp := make([]MemberResponse, len(members))
	for i, m := range members {
		u := userMap[m.UserID]
		resp[i] = MemberResponse{
			ID:       m.ID,
			UserID:   strconv.FormatInt(m.UserID, 10),
			Email:    u.Email,
			Name:     u.Name,
			Role:     m.Role,
			JoinedAt: m.JoinedAt.UTC().Format(time.RFC3339),
		}
	}

	return resp, nil
}

func (s *Service) memberRole(ctx context.Context, spaceID string, userID int64) (string, error) {
	var member schemas.SpaceMember
	if err := s.orm.WithContext(ctx).Where("space_id = ? AND user_id = ?", spaceID, userID).First(&member).Error; err != nil {
		return "", errors.Forbidden("not a member of this space")
	}
	return member.Role, nil
}

func toSpaceResponse(space schemas.Space, role string) SpaceResponse {
	return SpaceResponse{
		ID:          space.ID,
		Name:        space.Name,
		Description: space.Description,
		Role:        role,
		CreatedAt:   space.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   space.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func toSpaceDetailResponse(space schemas.Space, role string, members []MemberResponse) SpaceResponse {
	resp := toSpaceResponse(space, role)
	resp.Members = members
	return resp
}

func parseUserID(raw string) (int64, error) {
	uid, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, errors.Invalid("invalid user id")
	}
	return uid, nil
}
