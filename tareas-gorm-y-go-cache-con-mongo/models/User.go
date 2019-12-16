package models

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

/* ============================================================ USERS === */

// Credentials for saving authentication/authorization data
type Credentials struct {
	Username string `gorm:"unique;not null" json:"username" valid:"required, length(1|15)"`
	Password string `json:"password" valid:"required, length(1|20)"`
}

// User structure for user profile
type User struct {
	gorm.Model
	Credentials
	Fullname	string	`json:"fullname"`
	AvatarID	uint	`json:"avatar_id"gorm:"not null; type:int REFERENCES image(image_id) ON DELETE CASCADE"`
	NumPictures uint	`json:"num_pics"`
	NumComments uint	`json:"num_comments"`
	NumLikes	uint	`json:"num_likes"`
	Gallery	*[]Gallery 	`json:"gallery" gorm:"foreignkey:ImageID"`
	Images 		Image   `gorm:"foreignkey:AvatarID"`
}

// Checking if user has a Valid structure
func (u User) ValidUser() bool {
	ok, err := govalidator.ValidateStruct(u)
	fmt.Printf("Data entry (username and password) are both valid: %v /n", ok)
	return err == nil && ok
}

// NewCredentialsJSON for converting JSON content from the body to Credential object
func NewCredentialsJSON(jsonBytes []byte) *Credentials {
	cred := new(Credentials)
	err := json.Unmarshal(jsonBytes, cred)
	if err == nil {
		return cred
	}
	return nil
}

// NewUserJSON for converting JSON content from the body to User object
func NewUserJSON(jsonBytes []byte) *User {
	user := new(User)
	err := json.Unmarshal(jsonBytes, user)
	if err == nil {
		return user
	}
	return nil
}

// AddUser adds new user info to the Data Base (includes a new line to the table)
func AddUser(newUser *User, db *gorm.DB) (err error) {
	err = db.Create(newUser).Error
	return err
}

// GetUser looks for a specific username in Data Base
func GetUser(username string, db *gorm.DB) *User {
	user := new(User)
	db.Find(user, username)
	if user.Username == username {
		fmt.Println("User belongs to the Data Base")
		return user
	}
	fmt.Println("User does NOT belong to the Data Base")
	return nil
}