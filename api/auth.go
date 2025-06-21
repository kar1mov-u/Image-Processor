package api

import (
	"context"
	"encoding/json"
	"fmt"
	"image-processor/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (api *API) signIn(w http.ResponseWriter, r *http.Request) {

	//get the user details {Email, username, password}
	var userInfo models.SignIn
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

	//hash password
	userInfo.Password, err = HashPass(userInfo.Password)
	if err != nil {
		respondWithErr(w, r, http.StatusBadRequest, "password cannot be empty")
		return
	}

	//add to DB
	err = api.db.InsertUser(&userInfo)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			respondWithErr(w, r, http.StatusBadRequest, "Email is already in use")
			return
		} else {
			respondWithErr(w, r, http.StatusInternalServerError, "error on craeing user")
			return
		}
	}
	respondWithJson(w, r, 200, map[string]string{"message": "User created successfully", "user_id": userInfo.ID.String()})
}

func (api *API) login(w http.ResponseWriter, r *http.Request) {
	//user sends {Email, password(encrypted)}
	var loginInfo models.UserLogin

	err := json.NewDecoder(r.Body).Decode(&loginInfo)
	if err != nil {
		api.log.WithField("component", "api").Error(fmt.Sprintf("Failed to unmarshall user response %v", err))
		respondWithErr(w, r, http.StatusInternalServerError, "failed to unmarshall user response")
		return
	}

	//validate email
	if !isValidMail(loginInfo.Email) {
		respondWithErr(w, r, http.StatusBadRequest, "not valid email")
		return
	}

	//get data from DB
	userDB, err := api.db.GetUserByEmail(loginInfo.Email)
	if err != nil {
		respondWithErr(w, r, 404, err.Error())
		return
	}

	//check password
	if !ValidatePass(loginInfo.Password, userDB.Password) {
		respondWithErr(w, r, 403, "invalid email or password")
		return
	}
	//create JWT and return it
	claims := jwt.MapClaims{
		"sub": userDB.ID.String(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	singnedToken, err := token.SignedString([]byte(api.jwtKey))
	if err != nil {
		respondWithErr(w, r, http.StatusInternalServerError, "Error on generating access token")
		return
	}
	respondWithJson(w, r, 200, map[string]string{"token": singnedToken})
}

func JwtMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing token", 400)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid token format", 400)
				return
			}
			tokenStr := parts[1]

			token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "invalid token payload", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
