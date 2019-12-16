package data

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/models"
)
// ConnectDB creates a URI with parameters and associates it to postgres in order to open a specific data base
func ConnectDB() (*gorm.DB, error) {
	// The following arguments are ordered as follows: dbHost, userName, dbName and password
	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", "localhost", "social", "smimages", "media")
	return gorm.Open("postgres", dbUri)
}

// InitDB manages a data base connection and create tables
func InitDB() {
	dbCnx, err := ConnectDB()
	if err == nil {
		defer dbCnx.Close()
		dbCnx.AutoMigrate(&models.User{}, &models.Image{}, &models.Picture{}, &models.Comment{})
	} else {
		panic(err)
	}
}