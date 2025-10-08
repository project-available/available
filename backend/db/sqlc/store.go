package db

import "database/sql"

type Store interface {
	Querier
}
type SQLStore struct {
	*Queries
}
type SQLStoreq1 struct {
	*Queries
}
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
	}
}

