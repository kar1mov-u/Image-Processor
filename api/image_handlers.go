package api

import (
	"image-processor/internal/models"
	"io"
	"net/http"
	"os"
)

func (api *API) postImage(w http.ResponseWriter, r *http.Request) {

	userID, err := getUserID(r)
	if err != nil {
		respondWithErr(w, r, http.StatusBadRequest, "invalid userID")
		return
	}

	//image is send via multipart-data
	//parse the mulitpart data, so we can access
	r.ParseMultipartForm(10 << 20)

	// access the file from request
	//file is the io.Reader object(actual file) and header contains info about the file
	file, header, err := r.FormFile("image")
	if err != nil {
		api.log.WithField("component", "api").Errorf("failure in accessing image: %w", err)
		respondWithErr(w, r, 400, "error on uploading file")
		return
	}
	defer file.Close()

	//create the new location
	dest, err := os.Create("./uploads/" + header.Filename)
	if err != nil {
		api.log.WithField("component", "api").Errorf("failure in creaing image copy: %v", err)
		respondWithErr(w, r, http.StatusInternalServerError, "error on creaitng file copy")
		return
	}
	//copy file to the destination
	_, err = io.Copy(dest, file)
	if err != nil {
		api.log.WithField("component", "api").Errorf("failure on copying image : %v", err)
		respondWithErr(w, r, http.StatusInternalServerError, "error on copying file")
		return
	}
	defer dest.Close()

	//save into DB
	image := models.Image{Name: header.Filename, UserID: userID}
	err = api.db.InsertImage(&image)
	if err != nil {
		respondWithErr(w, r, http.StatusInternalServerError, "fail to insert image to db")
		return
	}

	respondWithJson(w, r, 200, map[string]any{"success": true, "image": image})
}

func (api *API) getImages(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		respondWithErr(w, r, http.StatusBadRequest, "invalid userID")
		return
	}

	images, err := api.db.GetImagesByUserID(userID)
	if err != nil {
		respondWithErr(w, r, 500, "fail to retrieve images")
		return
	}
	respondWithJson(w, r, 200, images)
}
