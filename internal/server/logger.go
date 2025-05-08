package server

import (
	"os"
	"time"
	"user-management/internal/config"

	log "github.com/sirupsen/logrus"
)

func SetupLogger() {
	logDir := config.LogPath
	logFile := logDir + "/app.log"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			log.Fatal("Failed to create log directory: ", err)
		}
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}

	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "time",
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "msg",
		},
	})
	log.SetLevel(log.InfoLevel)       
	log.SetOutput(file)

	log.WithFields(log.Fields{
		"service": "user-management",
		"env":     "development",
	}).Info("Logger initialized")
}