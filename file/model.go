package file

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	//import sqlite
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
)

// Source represents a source file that we want to follow changes
type Source struct {
	ID          string
	Path        string
	Fingerprint string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DocID       string
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
	FILE = iota
	URL
)

// Config represents a config line
type Config struct {
	Doc     Doc
	Sources []Source
}

// Result represents what we get after comparison between database and actual files
type Result struct {
	Doc     Doc
	Sources []Source
	Status  map[SourceStatus][]string
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

// NewSources create new sources recording file fingerprint
func NewSources(doc *Doc, sourcePaths []string) *[]Source {
	sources := []Source{}

	for _, path := range sourcePaths {
		fingerprint, err := calculateFingerprint(path)

		if err != nil {
			logrus.Fatal(err)
		}

		source := Source{
			uuid.NewV4().String(),
			path,
			fingerprint,
			time.Now(),
			time.Now(),
			doc.ID,
		}

		sources = append(sources, source)
	}

	return &sources
}

// InsertConfig create a new config entry
func InsertConfig(doc *Doc, sources *[]Source) {
	_, err := db.Exec("insert into docs values (?,?,?,?)", doc.ID, doc.Category, doc.Identifier, doc.CreatedAt)

	if err != nil {
		logrus.Fatal(err)
	}

	for _, source := range *sources {
		_, err := db.Exec("insert into sources values (?,?,?,?,?,?)", source.ID, source.Path, source.Fingerprint, source.CreatedAt, source.UpdatedAt, doc.ID)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}

// ListConfig return a config list
func ListConfig() *[]Config {
	configs := []Config{}

	rows, err := db.Query("select d.id, d.category, d.identifier, d.created_at, s.id, s.path, s.fingerprint, s.created_at, s.updated_at, s.doc_id from docs d inner join sources s on s.doc_id = d.id order by d.created_at")

	if err != nil {
		logrus.Fatal(err)
	}

	for rows.Next() {
		doc := Doc{}
		source := Source{}

		err := rows.Scan(&doc.ID, &doc.Category, &doc.Identifier, &doc.CreatedAt, &source.ID, &source.Path, &source.Fingerprint, &source.CreatedAt, &source.UpdatedAt, &source.DocID)

		if err != nil {
			logrus.Fatal(err)
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
			logrus.Fatal(err)
		}

		_, err = db.Exec(fmt.Sprintf("delete from docs where id in (%s)", strings.Join(docIds, ",")))

		if err != nil {
			logrus.Fatal(err)
		}
	}
}
