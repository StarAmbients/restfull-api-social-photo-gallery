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
	"strconv"
)

/* ========================================================= COMMENTS === */

// CreateComment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	comm := models.NewCommentJSON(jsonBytes)
	if comm == nil {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	if comm.ValidComment() == false {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}

	// We open the DB connection
	db, _ := data.ConnectDB()
	defer db.Close()

	// We add a new line to the table Image with all values in the structure of image
	if err := models.AddComment(comm, db); err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		println(fmt.Sprintf("Error creating comment: %s", err))
		return
	}

	// At this point there is no problem with the request
	// Writing the Header of the message (setting new route and specifying content type)
	// TODO We need to concatenate strings here "gallery/images/{id:[0-9]+} plus what follows as format
	w.Header().Set("Location", fmt.Sprintf("/comments/%d", comm.ID))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201
}

// GetComments
func GetComments(w http.ResponseWriter, r *http.Request) {
	db, _ := data.ConnectDB()
	defer db.Close()
	comm := models.GetComments(db)
	if comm != nil {
		jsonComm, err := json.Marshal(comm)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonComm)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

//GetComment
func GetComment(w http.ResponseWriter, r *http.Request) {
	if idStr, ok := mux.Vars(r)["comm_id"]; ok {
		db, _ := data.ConnectDB()
		defer db.Close()
		id, _ := strconv.Atoi(idStr)
		comm := models.GetComment(uint(id), db)
		if comm != nil {
			jsonComment, err := json.Marshal(comm)
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonComment)
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

//EditComment
func EditComment(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["comm_id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			comm := models.GetComment(uint(id), db)
			if comm != nil {
				if comm.UserID == userValid.(*models.User).ID {
					jsonBytes, err := ioutil.ReadAll(r.Body)
					if err == nil {
						editComm := new(models.Comment)
						err := json.Unmarshal(jsonBytes, editComm)
						if err == nil && comm.ValidComment() {
							editComm.UserID = userValid.(*models.User).ID
							editComm.ID = comm.ID
							models.EditComment(editComm, db)
							w.WriteHeader(http.StatusNoContent)
						} else {
							w.WriteHeader(http.StatusBadRequest)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
					}
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

//DeleteComment
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if userValid := r.Context().Value(middlewares.UserKey); userValid != nil {
		if idStr, ok := mux.Vars(r)["comm_id"]; ok {
			id, _ := strconv.Atoi(idStr)
			db, _ := data.ConnectDB()
			defer db.Close()
			comm := models.GetComment(uint(id), db)
			if comm != nil {
				if comm.UserID == userValid.(*models.User).ID {
					jsonBytes, err := ioutil.ReadAll(r.Body)
					if err == nil {
						deleteComm := new(models.Comment)
						err := json.Unmarshal(jsonBytes, deleteComm)
						if err == nil && comm.ValidComment() {
							deleteComm.UserID = userValid.(*models.User).ID
							deleteComm.ID = comm.ID
							models.EditComment(deleteComm, db)
							w.WriteHeader(http.StatusNoContent)
						} else {
							w.WriteHeader(http.StatusBadRequest)
						}
					} else {
						w.WriteHeader(http.StatusBadRequest)
					}
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

