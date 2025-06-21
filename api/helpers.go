package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func HashPass(password string) (string, error) {
	if len(password) == 0 {
		return "", fmt.Errorf("Password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePass(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserID(r *http.Request) (uuid.UUID, error) {
	userIdStr := r.Context().Value("userID").(string)
	userID, err := uuid.Parse(userIdStr)
	return userID, err
}
