package database

import (
	"fmt"
	"image-processor/internal/models"

	"github.com/google/uuid"
)

func (db *Database) InsertUser(userData *models.SignIn) error {
	query := `
	INSERT INTO users (username, email, password) VALUES (:username, :email, :password) RETURNING id
	`
	rows, err := db.conn.NamedQuery(query, userData)
	if err != nil {
		db.log.WithField("component", "db").Error(fmt.Sprintf("failed to insert user : %v", err))
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&userData.ID)
		if err != nil {
			db.log.WithField("component", "db").Errorf("failed to get user insert id: %v", err)
			return err
		}
	}

	return nil
}

func (db *Database) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := db.conn.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return user, fmt.Errorf("Invalid email or password")
	}
	return user, nil
}

func (db *Database) GetUserByID(userID uuid.UUID) (models.User, error) {
	var user models.User
	err := db.conn.Get(&user, "SELECT * FROM users WHERE id=$1", userID)
	if err != nil {
		return user, fmt.Errorf("User Not Found :%w", err)
	}
	return user, nil
}
