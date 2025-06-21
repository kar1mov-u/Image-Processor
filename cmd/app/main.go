package main

import (
	"fmt"
	"image-processor/api"
	"image-processor/internal/database"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()
	jwtKey := os.Getenv("JWT_KEY")
	log := log.New()

	//connect to the DB
	dsn := "host=localhost port=5433 user=appuser password=secret dbname=image_service sslmode=disable"
	dbConn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.WithField("component", "DB").Fatal(fmt.Sprintf("Failed to connect DB : %v", err))
	}
	log.WithField("component", "DB").Info("Connected to DB")
	defer dbConn.Close()

	db := database.New(dbConn, log)
	api := api.New(log, db, jwtKey)

	log.WithField("component", "api").Info("Starting Server..")
	http.ListenAndServe(":9999", api.Routes)
}
