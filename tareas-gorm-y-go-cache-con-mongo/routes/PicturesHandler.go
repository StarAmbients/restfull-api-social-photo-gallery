package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/data"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

/* ========================================================= PICTURES === */

// GetPictures
func GetPictures(w http.ResponseWriter, r *http.Request) {
	db, _ := data.ConnectDB()
	defer db.Close()
	pict := models.GetPictures(db)
	if pict != nil {
		jsonPict, err := json.Marshal(pict)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonPict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// CreatePicture
func CreatePicture(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	picture := models.NewPictureJSON(jsonBytes)
	if picture == nil {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	if picture.ValidPicture() == false {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}

	// We open the DB connection
	db, _ := data.ConnectDB()
	defer db.Close()

	// We add a new line to the table Image with all values in the structure of image
	if err := models.AddPicture(picture, db); err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		println(fmt.Sprintf("Error creating picture: %s", err))
		return
	}

	// At this point there is no problem with the request
	// Writing the Header of the message (setting new route and specifying content type)
	w.Header().Set("Location", fmt.Sprintf("/pictures/%d", picture.ID))
	w.Header().Add("Content-Type", "application/json")
	// Writing the body of the message (token) DSR not sure
	//w.Write(responseBytes) tibor
	w.WriteHeader(http.StatusCreated) //201
}

//GetPicture
func GetPicture(w http.ResponseWriter, r *http.Request) {
	if idStr, ok := mux.Vars(r)["id"]; ok {
		db, _ := data.ConnectDB()
		defer db.Close()
		id, _ := strconv.Atoi(idStr)
		pict := models.GetPicture(uint(id), db)
		if pict != nil {
			jsonPicture, err := json.Marshal(pict)
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonPicture)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}