package models

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

/* =========================================================== IMAGES === */

type Image struct {
	gorm.Model

	UserID			uint    	`json:"user_id" gorm:"not null, type:int REFERENCES user(user_id) ON DELETE CASCADE"`
	ThumbUUID   	uuid.UUID 	`json:"-"`
	LowResUUID  	uuid.UUID 	`json:"-"`
	HighResUUID 	uuid.UUID 	`json:"-"`
	ThumbURL    	string    	`gorm:"-" json:"thumb_url"`
	LowResURL   	string    	`gorm:"-" json:"lowres_url"`
	HighResURL  	string    	`gorm:"-" json:"highres_url"`
	User			*User 		`json:"user, omitempty" "gorm:foreignkey:UserID" `
}

// Checking if image has a Valid structure
func (i Image) ValidImage() bool {
	ok, err := govalidator.ValidateStruct(i)
	return err == nil && ok
}

// NewImageJSON for converting JSON content from the body to Image object
func NewImageJSON(jsonBytes []byte) *Image {
	image := new(Image)
	err := json.Unmarshal(jsonBytes, image)
	if err == nil {
		return image
	}
	return nil
}

// NewImage creates a new image
func NewImage(userID uint, thumb string, lowRes string, highRes string, thumbUUID uuid.UUID, LowResUUID uuid.UUID, HighResUUID uuid.UUID) *Image {
	return &Image{
		UserID: userID,
		ThumbURL: thumb,
		LowResURL: lowRes,
		HighResURL: highRes,
		ThumbUUID: thumbUUID,
		LowResUUID: LowResUUID,
		HighResUUID: HighResUUID,
	}
}

// AddImage adds new image info to the Data Base (includes a new line to the table)
func AddImage(newImage *Image, db *gorm.DB) (err error) {
	err = db.Create(newImage).Error
	return err
}

//GetImage to look for an image in the data base
func GetImage(id uint, db *gorm.DB) *Image {
	image := new(Image)
	db.Find(image, id)
	if image.ID == id {
		return image
	}
	return nil
}

// AfterFind is a suggestion from Fran
func (u *Image) AfterFind() (err error) {
	u.ThumbURL = fmt.Sprintf("/images/%s.jpg", u.ThumbUUID.String())
	u.LowResURL = fmt.Sprintf("/images/%s.jpg", u.LowResUUID.String())
	u.HighResURL = fmt.Sprintf("/images/%s.jpg", u.HighResUUID.String())
	return
}
