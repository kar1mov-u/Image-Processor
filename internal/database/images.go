package database

import (
	"image-processor/internal/models"

	"github.com/google/uuid"
)

func (db *Database) InsertImage(image *models.Image) error {
	query := `INSERT INTO images (user_id, name, uploaded_at, updated_at) VALUES (:user_id, :name,now(), now()) RETURNING id`

	rows, err := db.conn.NamedQuery(query, image)
	if err != nil {
		db.log.WithField("component", "db").Errorf("fail to insert image: %v", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&image.ID)
		if err != nil {
			db.log.WithField("component", "db").Errorf("fail to get imageID in inserting: %v", err)
			return err
		}
	}
	return nil
}

func (db *Database) GetImagesByUserID(userID uuid.UUID) ([]models.Image, error) {
	images := []models.Image{}

	query := `SELECT * FROM images WHERE user_id=$1`
	rows, err := db.conn.Query(query, userID)
	if err != nil {
		db.log.WithField("component", "db").Errorf("fail to query images %v", err)
		return images, nil
	}
	defer rows.Close()

	for rows.Next() {
		var image models.Image
		err = rows.Scan(&image.ID, &image.UserID, &image.Name, &image.Location, &image.UploadedAt, &image.UpdatedAt)
		if err != nil {
			db.log.WithField("component", "db").Errorf("fail to scan images while querying: %v", err)
			return images, err
		}
		images = append(images, image)
	}
	return images, nil
}
