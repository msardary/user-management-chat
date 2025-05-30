// @title User Management API
// @version 1.0
// @description This is a user management service and chat
// @host 91.107.176.236:8085
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"user-management/internal/config"
	"user-management/internal/db"
	"user-management/internal/server"
	dbg "user-management/internal/db/generated"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.LoadConfig()

	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	queries := dbg.New(dbConn)

	if err := server.Start(queries); err != nil {
		log.Fatal("Failed to start server: ", err)	
	}
}
