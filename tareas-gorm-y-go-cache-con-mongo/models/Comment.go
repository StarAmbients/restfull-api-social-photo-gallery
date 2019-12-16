package models

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
	"time"
)

/* ========================================================= COMMENTS === */

type Comment struct {
	gorm.Model
	UserID         	uint		`json:"user_id"`
	PictureID      	uint        `json:"picture_id"`
	Created			time.Time	`json:"created"`
	Comment        	string      `json: "comment"`
	User			User 		`json:"user, omitempty" gorm:"foreignkey:UserID"`
	Picture			Picture		`json:"picture" gorm:"foreignkey:PictureID"`
}

// Checking if Comment has a Valid structure
func (c Comment) ValidComment() bool {
	ok, err := govalidator.ValidateStruct(c)
	return err == nil && ok
}

// newComment creates a new comment
func newComment(userID uint, pictureID uint, created time.Time, comment string) *Comment {
	return &Comment{
		UserID:         userID,
		PictureID:      pictureID,
		Created:		created,
		Comment:        comment,
	}
}

// NewCommentJSON for converting JSON content from the body to Comment object
func NewCommentJSON(jsonBytes []byte) *Comment {
	comm := new(Comment)
	err := json.Unmarshal(jsonBytes, comm)
	if err == nil {
		return comm
	}
	return nil
}

//GetComments looks for all comments ordered by creation in the DB
func GetComments(db *gorm.DB) []Comment {
	var comments []Comment
	db.Preload("Picture").Preload("User").Order("created_at desc").Find(&comments)
	return comments
}

// GetComment gives us a specific comment
func GetComment(id uint, db *gorm.DB) *Comment {
	comm := new(Comment)
	db.Find(comm, id)
	if comm.ID == id {
		return comm
	}
	return nil
}

// AddComment adds new comment info to the Data Base (includes a new line to the table)
func AddComment(newComment *Comment, db *gorm.DB) (err error) {
	err = db.Create(newComment).Error
	return err
}

// EditComment allows edition of a comment into the Data Base (note: only Admin allowed)
func EditComment(editComment *Comment, db *gorm.DB) (err error) {
	err = db.Save(editComment).Error
	return err
}

// DeleteComment deletes a comment from the Data Base (note: only Admin allowed)
func DeleteComment(delComment *Comment, db *gorm.DB) (err error) {
	db.Delete(delComment)
	return err
}