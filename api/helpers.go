package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
)

func respondWithJson(w http.ResponseWriter, r *http.Request, statusCode int, payload any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Failed to send JSON response to client :%v. | Error :%v\n", r.RemoteAddr, err)
		http.Error(w, fmt.Sprintf("Failed to marshall repsonse: %v", err), http.StatusInternalServerError)
		return
	}

}

func respondWithErr(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	respondWithJson(w, r, statusCode, map[string]string{"error": message})

}

func isValidMail(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
