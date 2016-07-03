package file

import (
	"fmt"

	"strings"
	"time"

	//import sqlite
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
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
	SFILE = iota
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
	DFILE = iota
	DURL
)

// Config represents a config line
type Config struct {
	Doc     Doc
	Sources []Source
}

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

// NewItems create new source
func NewItems(identifiers *[]string, source *Source) *[]Item {
	items := []Item{}

	for _, identifier := range *identifiers {
		fingerprint, err := calculateFingerprint(identifier)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
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

	return &items
}

// InsertDoc create a new doc entry
func InsertDoc(doc *Doc) {
	_, err := db.Exec("insert into docs values (?,?,?,?)", doc.ID, doc.Category, doc.Identifier, doc.CreatedAt)

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}
}

// InsertSource create a new source entry
func InsertSource(source *Source) {
	_, err := db.Exec("insert into sources values (?,?,?,?,?)", source.ID, source.Identifier, source.Category, source.CreatedAt, source.DocID)

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}
}

func InsertItems(items *[]Item) {
	for _, item := range *items {
		_, err := db.Exec("insert into items values (?,?,?,?,?,?)", item.ID, item.Identifier, item.Fingerprint, item.CreatedAt, item.UpdatedAt, item.SourceID)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}
	}
}

// ListConfig return a config list
func ListConfig() *[]Config {
	configs := []Config{}

	rows, err := db.Query("select d.id, d.category, d.identifier, d.created_at, s.id, s.category, s.identifier, s.created_at, s.doc_id from docs d inner join sources s on s.doc_id = d.id order by d.created_at")

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}

	for rows.Next() {
		doc := Doc{}
		source := Source{}

		err := rows.Scan(&doc.ID, &doc.Category, &doc.Identifier, &doc.CreatedAt, &source.ID, &source.Category, &source.Identifier, &source.CreatedAt, &source.DocID)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		if len(configs) == 0 || configs[len(configs)-1].Doc.ID != source.DocID {
			configs = append(configs, Config{
				Doc: doc,
			})
		}

		configs[len(configs)-1].Sources = append(configs[len(configs)-1].Sources, source)
	}

	return &configs
}

func getItems(source *Source) *[]Item {
	items := []Item{}

	rows, err := db.Query("select id, identifier, fingerprint, created_at, updated_at, source_id from items where source_id = ? order by identifier", source.ID)

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}

	for rows.Next() {
		item := Item{}

		err := rows.Scan(&item.ID, &item.Identifier, &item.Fingerprint, &item.CreatedAt, &item.UpdatedAt, &item.SourceID)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		items = append(items, item)
	}

	return &items
}

// RemoveConfigs delete one or several config group
func RemoveConfigs(configs *[]Config) {
	sourceIds := []string{}
	docIds := []string{}

	for _, config := range *configs {
		for _, source := range config.Sources {
			sourceIds = append(sourceIds, fmt.Sprintf(`"%s"`, source.ID))
		}
		docIds = append(docIds, fmt.Sprintf(`"%s"`, config.Doc.ID))
	}

	if len(sourceIds) > 0 {
		_, err := db.Exec(fmt.Sprintf("delete from sources where id in (%s)", strings.Join(sourceIds, ",")))

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		_, err = db.Exec(fmt.Sprintf("delete from docs where id in (%s)", strings.Join(docIds, ",")))

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		_, err = db.Exec(fmt.Sprintf("delete from items where source_id in (%s)", strings.Join(sourceIds, ",")))

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}
	}
}
