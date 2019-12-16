package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	image_manger "github.com/graux/image-manager"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/data"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/middlewares"
	"github.com/starambients/tareas-gorm-y-go-cache-con-mongo/models"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

/* =========================================================== IMAGES === */

// createImage
func CreateImage(w http.ResponseWriter, r *http.Request) {
	if imageValid := r.Context().Value(middlewares.UserKey); imageValid != nil {
		// Reading the body of the request (images comes in an array of bytes)
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}

		// Checking the context, if this image creation is for the user avatar or a simple picture
		imgType := "picture"
		// TODO later I will implement the two option possibility that is in the context
		// imgType := ctx.GetString("type")
		// if !common.GetStringInArrayVariadic(*imgType, "picture", "avatar") {
		//    w.WriteHeader(http.StatusBadRequest) //400
		//    return
		//}

		// Capturing the path of the petition
		imagesPath, err := filepath.Abs("./public/images")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
			// There is no image folder
			return
		}

		// NewImageJSON for converting JSON content from the body to Image object
		imgManager := image_manger.NewImageManager(imagesPath)

		var uuids []uuid.UUID
		if imgType == "picture"{
			uuids, err = imgManager.ProcessImageAs16by9(jsonBytes)
		}else {
			if imgType == "avatar"{
				uuids, err = imgManager.ProcessImageAsSquare(jsonBytes)
			}
		}
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError) //500
			// There is an error while processing the image
		}else {
			if uuids == nil || len(uuids) != 3 {
				w.WriteHeader(http.StatusInternalServerError) //500
				// Invalid image processing
			}
		}
		newImg := &models.Image{
			//UserID: *ctx.UserId,		----> estoy con dificuldad de conseguir poner aqui el userID...
			//UserID: User.ID,
			ThumbUUID: uuids[0],
			LowResUUID: uuids[1],
			HighResUUID: uuids[2],
		}

		// We open the DB connection
		db, _ := data.ConnectDB()
		defer db.Close()

		// We add a new line to the table Image with all values in the structure of image
		if err := models.AddImage(newImg, db); err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
			println(fmt.Sprintf("Error creating image: %s", err))
			return
		}

		// At this point there is no problem with the request
		// Writing the Header of the message (setting new route and specifying content type)
		w.Header().Set("Location", fmt.Sprintf("/images%d", newImg.ID))
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) //201
		jsonImage, _ := json.Marshal(newImg)
		w.Write(jsonImage)
	}
}

// WriteImage encodes an image 'img' in jpeg format and writes it into ResponseWriter
func WriteImage(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("Unable to encode image")
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("Unable to write image")
	}
}

func GetImage(w http.ResponseWriter, r *http.Request){
	if idStr, ok := mux.Vars(r)["id"]; ok {
		db, _ := data.ConnectDB()
		defer db.Close()
		id, _ := strconv.Atoi(idStr)
		img := models.GetImage(uint(id), db)
		if img != nil {
			jsonImg, err := json.Marshal(img)
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonImg)
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
