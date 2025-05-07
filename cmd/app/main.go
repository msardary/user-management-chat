package main

import (
	"context"
	"user-management/internal/config"
	"user-management/internal/db"
	dbg "user-management/internal/db/generated"
	"user-management/internal/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	config.LoadConfig()

	dbConn, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	conn, err := dbConn.Acquire(context.Background())
	if err != nil {
		log.Fatal("Failed to acquire connection: ", err)
	}
	defer conn.Release()
	
	queries := dbg.New(conn)

	if err := server.Start(queries); err != nil {
		log.Fatal("Failed to start server: ", err)	
	}
}
