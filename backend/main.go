package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/project-available/available/api"
	db "github.com/project-available/available/db/sqlc"
	_ "github.com/project-available/available/docs"
	"github.com/project-available/available/utils"
)

// @title Available API
// @version 1.0
// @description API for Available - Room Booking Application for Students
// @description This API provides endpoints for managing room bookings, user accounts, and related resources.

// @contact.name API Support
// @contact.email support@available.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the access token.

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
