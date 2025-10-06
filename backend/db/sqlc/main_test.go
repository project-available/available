package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)
const (
	dbDriver = "postgres"
	dbSource = "postgres://postgres:12345@localhost:5000/available?sslmode=disable"
)

var testQuery *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver,dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQuery = New(conn)

	os.Exit(m.Run())

}
