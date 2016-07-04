package file

import (
	"time"

	//import sqlite
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
)

// Source represents a source that we want to follow changes
type Source struct {
	ID         string
	Category   SourceCategory
	Identifier string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DocID      string
}

// SourceCategory represents a source category
type SourceCategory int

// SourceCategory categories
const (
	SERROR = iota
	SFILE
	SFOLDER
)

// Item represents an actual tracked source
type Item struct {
	ID          string
	Identifier  string
	Fingerprint string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SourceID    string
}

// Doc represents a document which as relationship with one or several source files
type Doc struct {
	ID         string
	Category   DocCategory
	Identifier string
	CreatedAt  time.Time
}

// DocCategory represents a doc category
type DocCategory int

// DocCategory categories
const (
	DERROR = iota
	DFILE
	DURL
)

// Result represents what we get after comparison between database and actual files
type Result struct {
	Doc    Doc
	Source Source
	Status map[ItemStatus][]string
}

// NewDoc create a new doc file
func NewDoc(identifier string, category DocCategory) *Doc {
	return &Doc{
		uuid.NewV4().String(),
		category,
		identifier,
		time.Now(),
	}
}

// NewSource create new source
func NewSource(doc *Doc, identifier string, category SourceCategory) *Source {
	source := Source{
		uuid.NewV4().String(),
		category,
		identifier,
		time.Now(),
		time.Now(),
		doc.ID,
	}

	return &source
}

// NewItems create several new items
func NewItems(identifiers *[]string, source *Source) (*[]Item, error) {
	items := []Item{}

	for _, identifier := range *identifiers {
		fingerprint, err := calculateFingerprint(identifier)

		if err != nil {
			return nil, err
		}

		items = append(items, Item{
			uuid.NewV4().String(),
			identifier,
			fingerprint,
			time.Now(),
			time.Now(),
			source.ID,
		})
	}

	return &items, nil
}

// InsertDoc create a new doc entry
func InsertDoc(doc *Doc) error {
	_, err := db.Exec("insert into docs values (?,?,?,?)", doc.ID, doc.Category, doc.Identifier, doc.CreatedAt)

	return err
}

// InsertSource create a new source entry
func InsertSource(source *Source) error {
	_, err := db.Exec("insert into sources values (?,?,?,?,?)", source.ID, source.Identifier, source.Category, source.CreatedAt, source.DocID)

	return err
}

// InsertItems create severan new items
func InsertItems(items *[]Item) error {
	for _, item := range *items {
		_, err := db.Exec("insert into items values (?,?,?,?,?,?)", item.ID, item.Identifier, item.Fingerprint, item.CreatedAt, item.UpdatedAt, item.SourceID)

		if err != nil {
			return err
		}
	}

	return nil
}

func getItems(source *Source) (*[]Item, error) {
	items := []Item{}

	rows, err := db.Query("select id, identifier, fingerprint, created_at, updated_at, source_id from items where source_id = ? order by identifier", source.ID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := Item{}

		err := rows.Scan(&item.ID, &item.Identifier, &item.Fingerprint, &item.CreatedAt, &item.UpdatedAt, &item.SourceID)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return &items, nil
}
