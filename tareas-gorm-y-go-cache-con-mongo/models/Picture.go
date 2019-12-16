package models

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"time"
)

/* ========================================================= PICTURES === */

type Picture struct {
	gorm.Model

	UserID			uint    	`json:"user_id" gorm:"not null; type:int REFERENCES user(user_id) ON DELETE CASCADE"`
	ImageID			uint    	`json:"image_id" gorm:"not null; type:int REFERENCES image(image_id) ON DELETE CASCADE"`
	Title			string 		`json:"title" gorm:"not null; type:varchar(200)"`
	Description		*string 	`json:"description"`
	Created			time.Time	`json:"created" gorm:"not null; default:CURRENT_TIMESTAMP"`
	NumLikes		*uint		`json:"num_likes, omitempty"`
	NumComments		*uint		`json:"num_comments, omitempty"`
	Image			Image 		`json:"image" gorm:"foreignkey:ImageID"`
	User			*User 		`json:"user, omitempty" gorm:"foreignkey:UserID"`
}

// Checking if Picture has a Valid structure
func (p Picture) ValidPicture() bool {
	ok, err := govalidator.ValidateStruct(p)
	return err == nil && ok
}

// newPicture creates a new picture
func newPicture(userID uint, imageID uint, title string, description *string, created time.Time, numLikes *uint, numComments *uint) *Picture {
	return &Picture{
		UserID:			userID,
		ImageID:		imageID,
		Title:			title,
		Description:	description,
		Created:		created,
		NumLikes:		numLikes,
		NumComments:	numComments,
	}
}

// NewPictureJSON for converting JSON content from the body to Picture object
func NewPictureJSON(jsonBytes []byte) *Picture {
	picture := new(Picture)
	err := json.Unmarshal(jsonBytes, picture)
	if err == nil {
		return picture
	}
	return nil
}

// GetPictures looks for all images ordered by creation in the DB, and populates an array of pictures
func GetPictures(db *gorm.DB) []Picture {
	var pictures []Picture
	db.Preload("Image").Preload("User").Order("created_at desc").Find(&pictures)
	return pictures
}

// Given an id GetPicture looks for it specifically in the data base
func GetPicture(id uint, db *gorm.DB) *Picture {
	//Preloaded needed to give us the image and its id
	picture := new(Picture)
	db.Find(picture, id)
	if picture.ImageID == id {
		return picture
	}
	return nil
}

// AddPicture adds new picture info to the Data Base (includes a new line to the table)
func AddPicture(newPicture *Picture, db *gorm.DB) (err error) {
	err = db.Create(newPicture).Error
	return err
}