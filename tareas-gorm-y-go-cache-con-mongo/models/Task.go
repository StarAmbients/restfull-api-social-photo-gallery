package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

type Task struct {
	TaskID    int       `gorm:"primary_key" json:"task_id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title" valid:"required,length(3|50)"`
	Content   *string   `json:"content" valid:"length(3|500)"`
	Date      time.Time `json:"date"`
	Completed bool      `json:"completed"`
}

func (t Task) Valid() bool {
	ok, err := govalidator.ValidateStruct(t)
	return err == nil && ok
}

func NewTask(taskID int, userID int, title string, content string, date time.Time, completed bool) *Task {
	return &Task{
		TaskID:    taskID,
		UserID:    userID,
		Title:     title,
		Content:   &content,
		Date:      date,
		Completed: completed,
	}
}

func timeMustParse(dateStr string) time.Time {
	miTime, _ := time.Parse("2006-01-02", dateStr)
	return miTime
}

func GetTasks(db *gorm.DB) []Task {
	var tasks []Task
	db.Find(&tasks)
	return tasks
}

func GetTask(id int, db *gorm.DB) *Task {
	task := new(Task)
	db.Find(task, id)
	if task.TaskID == id {
		return task
	}
	return nil
}

func AddTask(newTask *Task, db *gorm.DB) {
	db.Create(newTask)
}

func EditTask(editTask *Task, db *gorm.DB) {
	db.Save(editTask)
}

func DeleteTask(delTask *Task, db *gorm.DB) {
	db.Delete(delTask)
}
