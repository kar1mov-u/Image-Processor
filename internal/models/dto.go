package models

import "github.com/google/uuid"

type UserLogin struct {
	Email    string
	Password string
}

type SignIn struct {
	Email    string
	Password string
	UserName string
	ID       uuid.UUID
}
