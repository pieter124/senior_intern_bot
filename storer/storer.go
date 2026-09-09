package storer

import "senior_intern_bot/domain"

type Storer interface {
	Store([]domain.Posting) error
}

type SQLiteStorer struct {
	// dbConnection sqlite3.dbconnection
}

// Pass in DB connection
func New() *SQLiteStorer {
	// TODO
	// Make connection to SQLite db here
	return &SQLiteStorer{
		// dbConnection: dbConnection
	}
}

func (storer *SQLiteStorer) Store(sentPostings []domain.Posting) error {
	// TODO
	// for , _ := range sentPostings {
	// Store each posting
	// }
	return nil
}
