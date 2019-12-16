package routes

import (
	"fmt"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/data"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/models"
	"io/ioutil"
	"net/http"
)

/* ============================================================ LIKES === */

// CreateLike
func CreateLike(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	like := models.NewLikeJSON(jsonBytes)
	if like == nil {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	if like.ValidLike() == false {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}

	// We open the DB connection
	db, _ := data.ConnectDB()
	defer db.Close()

	// We add a new line to the table Image with all values in the structure of image
	if err := models.AddLike(like, db); err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		println(fmt.Sprintf("Error creating like: %s", err))
		return
	}

	// At this point there is no problem with the request
	// Writing the Header of the message (setting new route and specifying content type)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201
}
