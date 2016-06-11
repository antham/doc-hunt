package file

import (
	"fmt"
	"strings"
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

// Config represents a config line
type Config struct {
	DocFile     Doc
	SourceFiles []Source
}

// Result represents what we get after comparison between database and actual files
type Result struct {
	DocFile     Doc
	SourceFiles []Source
	Status      map[string]string
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

// ListConfig return a config list
func ListConfig() *[]Config {
	configs := []Config{}

	rows, err := db.Query("select d.id, d.path, d.created_at, s.id, s.path, s.fingerprint, s.created_at, s.updated_at, s.doc_file_id from doc_file d inner join source_file s on s.doc_file_id = d.id order by d.created_at")

	if err != nil {
		logrus.Fatal(err)
	}

	for rows.Next() {
		docFile := Doc{}
		sourceFile := Source{}

		rows.Scan(&docFile.ID, &docFile.Path, &docFile.CreatedAt, &sourceFile.ID, &sourceFile.Path, &sourceFile.Fingerprint, &sourceFile.CreatedAt, &sourceFile.UpdatedAt, &sourceFile.DocFileID)

		if len(configs) == 0 || configs[len(configs)-1].DocFile.ID != sourceFile.DocFileID {
			configs = append(configs, Config{
				DocFile: docFile,
			})
		}

		configs[len(configs)-1].SourceFiles = append(configs[len(configs)-1].SourceFiles, sourceFile)
	}

	return &configs
}

// RemoveConfigs delete one or several config group
func RemoveConfigs(configs *[]Config) {
	sourceIds := []string{}
	docIds := []string{}

	for _, config := range *configs {
		for _, source := range config.SourceFiles {
			sourceIds = append(sourceIds, fmt.Sprintf(`"%s"`, source.ID))
		}
		docIds = append(docIds, fmt.Sprintf(`"%s"`, config.DocFile.ID))
	}

	if len(sourceIds) > 0 {
		_, err := db.Exec(fmt.Sprintf("delete from source_file where id in (%s)", strings.Join(sourceIds, ",")))

		if err != nil {
			logrus.Fatal(err)
		}

		_, err = db.Exec(fmt.Sprintf("delete from doc_file where id in (%s)", strings.Join(docIds, ",")))

		if err != nil {
			logrus.Fatal(err)
		}
	}
}
