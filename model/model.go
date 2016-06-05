package model

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Sirupsen/logrus"
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
	DocFileID   string
}

// Doc represents a document file which as relationship with one or several source files
type Doc struct {
	ID        string
	Path      string
	CreatedAt time.Time
}

// NewDoc create a new doc file
func NewDoc(docPath string) *Doc {
	return &Doc{
		uuid.NewV4().String(),
		docPath,
		time.Now(),
	}
}

// NewSources create new sources recording file fingerprint
func NewSources(doc *Doc, sourcePaths []string) *[]Source {
	sources := []Source{}

	for _, path := range sourcePaths {
		var fingerprint string

		data, err := ioutil.ReadFile(path)

		if err == nil {
			fingerprint = fmt.Sprintf("%x", sha1.Sum(data))
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
	_, err := db.Exec("insert into doc_file values (?,?,?)", doc.ID, doc.Path, doc.CreatedAt)

	if err != nil {
		logrus.Fatal(err)
	}

	for _, source := range *sources {
		_, err := db.Exec("insert into source_file values (?,?,?,?,?,?)", source.ID, source.Path, source.Fingerprint, source.CreatedAt, source.UpdatedAt, doc.ID)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}
