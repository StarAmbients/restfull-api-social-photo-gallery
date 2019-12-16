package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/data"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/lib"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/models"
	"io/ioutil"
	"net/http"
)

/* ============================================================ USERS === */

// Register creates a new user. It awaits the object User.
func Register(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	user := models.NewUserJSON(jsonBytes)
	if user == nil {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	if user.ValidUser() == false {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}

	// We open the DB connection
	db, _ := data.ConnectDB()
	defer db.Close()

	// We check if the username is already in the data base
	u := models.GetUser(user.Username, db)
	if u != nil{
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}

	// We add a new line to the table User with all values in the structure of user
	if err := models.AddUser(user, db); err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		println(fmt.Sprintf("Error creating user: %s", err))
		return
	}

	// At this point there is no problem with the request
	// Writing the Header of the message (setting new route and specifying content type)
	w.Header().Set("Location", fmt.Sprintf("/users/%s", user.Username))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201
}

// login requires object with credentials
func Login(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cred := models.NewCredentialsJSON(jsonBytes)
	if cred == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db, _ := data.ConnectDB()
	defer db.Close()
	validUser := lib.ValidateCredent(cred, db)
	if validUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// token, err := lib.CreateJWT(validUser)
	token, err := lib.CreateToken(validUser, data.GetCacheClient())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBytes, err := json.Marshal(lib.JSONToken{Token: token})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //200
	w.Write(responseBytes) //Token inside body
}

// getUser recovers user data via its username
func GetUser(w http.ResponseWriter, r *http.Request) {
	if username, ok := mux.Vars(r)["username"]; ok {
		db, _ := data.ConnectDB()
		defer db.Close()
		//id, _ := strconv.Atoi(idStr)
		user := models.GetUser(username, db)
		if user != nil {
			jsonUser, err := json.Marshal(user)
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonUser)
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

