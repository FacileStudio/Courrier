package schemas

import "time"

// Space is a shared workspace whose members collaborate on accounts.
type Space struct {
	ID          string        `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string        `gorm:"column:name;not null"`
	Description string        `gorm:"column:description"`
	CreatedAt   time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	Members     []SpaceMember `gorm:"foreignKey:SpaceID"`
}

func (Space) TableName() string { return "spaces" }

// SpaceMember is a user's membership, and role, in a space.
type SpaceMember struct {
	ID       string    `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	SpaceID  string    `gorm:"column:space_id;not null;index"`
	UserID   int64     `gorm:"column:user_id;not null;index"`
	Role     string    `gorm:"column:role;not null;default:'member'"`
	JoinedAt time.Time `gorm:"column:joined_at;autoCreateTime"`
	Space    Space     `gorm:"foreignKey:SpaceID;constraint:OnDelete:CASCADE"`
	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (SpaceMember) TableName() string { return "space_members" }

const (
	SpaceRoleOwner  = "owner"
	SpaceRoleAdmin  = "admin"
	SpaceRoleMember = "member"
)
