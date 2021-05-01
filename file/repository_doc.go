package file

import (
	"database/sql"
)

// DocRepository handle all actions availables on docs table
type DocRepository struct {
	db *sql.DB
}

// NewDocRepository return a new DocRepository instance
func NewDocRepository(db *sql.DB) DocRepository {
	return DocRepository{db: db}
}

// Create insert a new doc entry
func (d DocRepository) Create(doc *Doc) error {
	_, err := d.db.Exec("insert into docs values (?,?,?,?)", doc.ID, doc.Category, doc.Identifier, doc.CreatedAt)

	return err
}
