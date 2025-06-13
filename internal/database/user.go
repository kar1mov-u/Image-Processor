package database

import (
	"fmt"
	"image-processor/internal/models"
)

func (db *Database) InsertUser(userData *models.User) error {
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

func (db *Database) GetUser(email string) models.User {

}
