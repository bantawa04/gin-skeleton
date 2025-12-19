package user

import (
	"time"

	"gin/internal/constant"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID               string                   `json:"id" gorm:"primaryKey;type:char(26)"`
	FirstName        *string                  `json:"first_name,omitempty" gorm:"type:varchar(255)"`
	LastName         *string                  `json:"last_name,omitempty" gorm:"type:varchar(255)"`
	Email            string                   `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	Password         string                   `json:"-" gorm:"type:varchar(255)"`
	Phone            *string                  `json:"phone,omitempty" gorm:"type:varchar(20)"`
	Province         *string                  `json:"province,omitempty" gorm:"type:varchar(100)"`
	District         *string                  `json:"district,omitempty" gorm:"type:varchar(100)"`
	City             *string                  `json:"city,omitempty" gorm:"type:varchar(100)"`
	Zip              *string                  `json:"zip,omitempty" gorm:"type:varchar(10)"`
	Country          *string                  `json:"country,omitempty" gorm:"type:varchar(100)"`
	Address          *string                  `json:"address,omitempty" gorm:"type:varchar(255)"`
	Type             constant.AccountTypeEnum `json:"type" gorm:"type:varchar(20);default:'user'"`
	SocialProvider   *string                  `json:"social_provider,omitempty" gorm:"type:varchar(50);default:'emailPassword'"`
	SocialProviderID *string                  `json:"social_provider_id,omitempty" gorm:"type:varchar(100)"`
	LastSignInAt     *time.Time               `json:"last_sign_in_at,omitempty"`
	Status           constant.UserStatusEnum  `json:"status" gorm:"type:varchar(20);default:'inactive'"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
	DeletedAt        gorm.DeletedAt           `json:"deleted_at,omitempty" gorm:"index"`
}

// BeforeCreate hook for generating ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		// Generate a new ULID
		id := ulid.Make()
		u.ID = id.String()
	}
	return nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	firstName := ""
	if u.FirstName != nil {
		firstName = *u.FirstName
	}

	lastName := ""
	if u.LastName != nil {
		lastName = *u.LastName
	}

	if firstName != "" && lastName != "" {
		return firstName + " " + lastName
	} else if firstName != "" {
		return firstName
	} else if lastName != "" {
		return lastName
	}
	return ""
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
