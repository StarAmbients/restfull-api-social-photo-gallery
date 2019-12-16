package models

import (
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	UserID			uint    	`json:"user_id" gorm:"not null; type:int REFERENCES user(user_id) ON DELETE CASCADE"`
	Role			string		`json:"user`
	User			*User 		`json:"user, omitempty" gorm:"foreignkey:UserID"`
}

// Checking if Like has a Valid structure
func (r Role) ValidRole() bool {
	ok, err := govalidator.ValidateStruct(r)
	return err == nil && ok
}

//Commenting newLike
func newRole(userID uint, role string) *Role {
	return &Role{
		UserID:    userID,
		Role:      role,
	}
}
