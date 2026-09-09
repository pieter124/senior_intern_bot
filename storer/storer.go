package storer

import (
	"database/sql"
	"log"
	"senior_intern_bot/domain"
)

type Storer interface {
	Store([]domain.Posting) error
}

type SQLiteStorer struct {
	dbConnection *sql.DB
}

// Pass in DB connection
func New(dbConnection *sql.DB) *SQLiteStorer {
	return &SQLiteStorer{
		dbConnection: dbConnection,
	}
}

func (storer *SQLiteStorer) Store(postings []domain.Posting) error {
	// Start transaction
	tx, err := storer.dbConnection.Begin()
	if err != nil {
		log.Println("error beginning db transaction: ", err)
		return err
	}
	var query string = `INSERT OR IGNORE INTO postings
(source, company, job_id, title, url, location, verdict, posted_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	for _, p := range postings {
		_, err := tx.Exec(query, p.Source, p.Company,
			p.JobID, p.Title, p.URL, p.Location, int(p.Verdict), p.PostedAt, p.UpdatedAt)
		if err != nil {
			log.Printf("error trying to store posting: %s|%s|%s\n", p.Source, p.Company, p.JobID)
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println("error trying to commit transaction: ", err)
		return err
	}
	return nil
}
