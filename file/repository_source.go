package file

import (
	"database/sql"
)

// SourceRepository handle all actions availables on sources table
type SourceRepository struct {
	db *sql.DB
}

// NewSourceRepository return a new SourceRepository instance
func NewSourceRepository(db *sql.DB) SourceRepository {
	return SourceRepository{db: db}
}

// Create insert a new source entry
func (s SourceRepository) Create(source *Source) error {
	_, err := s.db.Exec("insert into sources values (?,?,?,?,?)", source.ID, source.Identifier, source.Category, source.CreatedAt, source.DocID)

	return err
}
