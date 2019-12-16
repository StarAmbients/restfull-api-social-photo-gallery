package models

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"time"
)

/* ============================================================ LIKES === */

type Like struct {
	gorm.Model
	UserID			uint    	`json:"user_id" gorm:"not null; type:int REFERENCES user(user_id) ON DELETE CASCADE"`
	PictureID		uint    	`json:"picture_id" gorm:"not null; type:int REFERENCES picture(picture_id) ON DELETE CASCADE"`
	Liked  			time.Time	`json:"when_liked"`
	User			User 		`json:"user, omitempty" gorm:"foreignkey:UserID"`
	Picture			Picture		`json:"picture" gorm:"foreignkey:PictureID"`
}

// Checking if Like has a Valid structure
func (l Like) ValidLike() bool {
	ok, err := govalidator.ValidateStruct(l)
	return err == nil && ok
}

// newLike creates a new Like
func NewLike(userID uint, pictureID uint, whenLiked time.Time) *Like {
	return &Like{
		UserID:         userID,
		PictureID:      pictureID,
		Liked:			whenLiked,
	}
}

// NewLikeJSON for converting JSON content from the body to Like object
func NewLikeJSON(jsonBytes []byte) *Like {
	like := new(Like)
	err := json.Unmarshal(jsonBytes, like)
	if err == nil {
		return like
	}
	return nil
}

// AddLike adds new Like info to the Data Base (includes a new line to the table)
func AddLike(newLike *Like, db *gorm.DB) (err error) {
	err = db.Create(newLike).Error
	return err
}
