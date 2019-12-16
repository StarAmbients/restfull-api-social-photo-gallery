package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/data"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/middlewares"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/models"
	"io/ioutil"
	"net/http"
)

/* ======================================================== GALLERIES === */

// createGallery
func CreateGallery(w http.ResponseWriter, r *http.Request) {
	if galleryValid := r.Context().Value(middlewares.UserKey); galleryValid != nil {
		// Reading the body of the request (gallery comes in an array of bytes)
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}

		// NewGalleryJSON for converting JSON content from the body to Image object
		gallery := models.NewGalleryJSON(jsonBytes)
		if gallery == nil {
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}
		if gallery.ValidGallery() == false {
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}
		// We open the DB connection
		db, _ := data.ConnectDB()
		defer db.Close()

		// We add a new line to the table Image with all values in the structure of image
		if err := models.AddGallery(gallery, db); err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
			println(fmt.Sprintf("Error creating gallery: %s", err))
			return
		}

		// At this point there is no problem with the request
		// Writing the Header of the message (setting new route and specifying content type)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) //201
		jsonGallery, _ := json.Marshal(gallery)
		w.Write(jsonGallery)
	}
}


// getGallery recovers user data via its username
func GetGallery(w http.ResponseWriter, r *http.Request) {
	if username, ok := mux.Vars(r)["username"]; ok {
		db, _ := data.ConnectDB()
		defer db.Close()
		gallery := models.GetGallery(username, db)
		if gallery != nil {
			jsonGallery, err := json.Marshal(gallery)
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonGallery)
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

