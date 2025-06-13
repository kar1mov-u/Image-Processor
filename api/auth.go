package api

import (
	"encoding/json"
	"fmt"
	"image-processor/internal/models"
	"net/http"
)

func (api *API) signIn(w http.ResponseWriter, r *http.Request) {

	//get the user details {Email, username, password}
	var userInfo models.User
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		api.log.WithField("component", "api").Error(fmt.Sprintf("Failed to unmarshall user response %v", err))
		respondWithErr(w, r, http.StatusInternalServerError, "failed to unmarshall user response")
		return
	}

	//validate email
	if !isValidMail(userInfo.Email) {
		respondWithErr(w, r, http.StatusBadRequest, "not valid email")
		return
	}

	//add to DB
	err = api.db.InsertUser(&userInfo)
	if err != nil {
		respondWithErr(w, r, http.StatusInternalServerError, "error on creating user")
		return
	}
	//return OK
	respondWithJson(w, r, 204, userInfo.ID)

}

func (api *API) login(w http.ResponseWriter, r *http.Request) {
	//user sends {Email, password(encrypted)}

	//marshall into Struct
	//validate
	//get data from DB
	//check password
	//response with JWT

}
