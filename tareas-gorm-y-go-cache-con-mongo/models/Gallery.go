package models

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

/* ======================================================== GALLERIES === */


type Gallery struct {
	gorm.Model
	UserID			uint    	`json:"user_id" gorm:"not null; type:int REFERENCES user(user_id) ON DELETE CASCADE"`
	PicturesIDs		[]uint    	`json:"picture_id"`
	Pictures		[]Picture	`json:"pictures"`
	User			User 		`json:"user, omitempty" gorm:"foreignkey:UserID"`
}

// Checking if Gallery has a Valid structure
func (g Gallery) ValidGallery() bool {
	ok, err := govalidator.ValidateStruct(g)
	return err == nil && ok
}

// newLike creates a new Gallery
func NewGallery(userID uint, picturesIDs []uint, pictures []Picture) *Gallery {
	return &Gallery{
		UserID:         userID,
		PicturesIDs:    picturesIDs,
		Pictures:		pictures,
	}
}

// NewGalleryJSON for converting JSON content from the body to Like object
func NewGalleryJSON(jsonBytes []byte) *Gallery {
	gallery := new(Gallery)
	err := json.Unmarshal(jsonBytes, gallery)
	if err == nil {
		return gallery
	}
	return nil
}

// GetGallery
func GetGallery(username string, db *gorm.DB) *Picture {
	gallery := new(Picture)
	db.Find(gallery, username)
	//if gallery.(*models.User).Username == username {
	//	return gallery
	//}
	return nil
}

// AddLike adds new Gallery info to the Data Base (includes a new line to the table)
func AddGallery(newGallery *Gallery, db *gorm.DB) (err error) {
	err = db.Create(newGallery).Error
	return err
}
