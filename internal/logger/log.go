package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func New() *log.Logger {
	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(log.InfoLevel)
	return logger
}
