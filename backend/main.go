package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/getsentry/sentry-go"
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
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://3c28e31dba3b2ff88d5e5ff2a12e5e33@o4510268002664448.ingest.de.sentry.io/4510702348533840",
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Debug:            true,
		SendDefaultPII:   true,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	defer sentry.Flush(2 * time.Second)

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// Optimize connection pool
	// conn.SetMaxOpenConns(25)
	// conn.SetMaxIdleConns(25)
	// conn.SetConnMaxLifetime(5 * 60 * 1000000000) // 5 minutes

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
