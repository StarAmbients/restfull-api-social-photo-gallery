package main

import (
	"github.com/gorilla/mux"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/data"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/middlewares"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/routes"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

// main code
func main() {

	// Initiating Router
	router := mux.NewRouter().StrictSlash(true)

	// Handling routes
	/* ============================================================ USERS === */
	router.HandleFunc("/users", routes.Register).Methods(http.MethodPost)
	router.HandleFunc("/users", routes.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/users/login", routes.Login).Methods(http.MethodPost)
	router.HandleFunc("/users/{username}", routes.GetUser).Methods(http.MethodGet)

	/* =========================================================== IMAGES === */
	router.HandleFunc("/images", routes.CreateImage).Methods(http.MethodPost)
	router.HandleFunc("/images/{id:[0-9]+}", routes.GetImage).Methods(http.MethodGet)

	/* ========================================================= PICTURES === */
	router.HandleFunc("/pictures", routes.GetPictures).Methods(http.MethodGet)
	router.HandleFunc("/pictures", routes.CreatePicture).Methods(http.MethodPost)
	router.HandleFunc("/pictures/{id:[0-9]+}", routes.GetPicture).Methods(http.MethodGet)

	/* ========================================================= COMMENTS === */
	router.HandleFunc("gallery/images/{id:[0-9]+}/comments", routes.CreateComment).Methods(http.MethodPost)
	router.HandleFunc("gallery/images/{id:[0-9]+}/comments/{comm_id:[0-9]+}", routes.EditComment).Methods(http.MethodPut)
	router.HandleFunc("gallery/images/{id:[0-9]+}/comments/{comm_id:[0-9]+}", routes.GetComment).Methods(http.MethodGet)
	router.HandleFunc("gallery/images/{id:[0-9]+}/comments/{comm_id:[0-9]+}", routes.DeleteComment).Methods(http.MethodDelete)

	/* ============================================================ LIKES === */
	router.HandleFunc("pictures/{pic_id:[0-9]+}/likes/{user_id:[0-9]+}", routes.CreateLike).Methods(http.MethodPost)

	/* ======================================================== GALLERIES === */
	router.HandleFunc("/gallery/images", routes.CreateGallery).Methods(http.MethodPost)
	router.HandleFunc("/gallery/users/{username}", routes.GetGallery).Methods(http.MethodGet)

	// Authenticating user
	router.Use(middlewares.AuthUser)

	// Creating Negroni (idiomatic approach to web middleware)
	middle := negroni.Classic()

	// Registering middlewares (used between a ServeMux and our application handlers)
	middle.UseHandler(router)

	// Initiating Data Base
	data.InitDB()

	// Starting an HTTP server with a given address and handler
	err := http.ListenAndServe(":8080", middle)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}