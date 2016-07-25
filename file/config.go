package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/antham/doc-hunt/util"
)

// Config represents a config line
type Config struct {
	Doc     Doc
	Sources []Source
}

func createFolderSource(source *Source) error {
	err := InsertSource(source)

	if err != nil {
		return err
	}

	files, err := util.ExtractFolderFiles(source.Identifier)

	if err != nil {
		return err
	}

	items, err := NewItems(files, source)

	if err != nil {
		return err
	}

	err = InsertItems(items)

	return err
}

func createFileSource(source *Source) error {
	err := InsertSource(source)

	if err != nil {
		return err
	}

	items, err := NewItems(&[]string{source.Identifier}, source)

	if err != nil {
		return err
	}

	err = InsertItems(items)

	return err
}

// ParseIdentifier extract identifier and category from string
func ParseIdentifier(value string) (string, SourceCategory) {
	f, ferr := os.Stat(util.GetAbsPath(value))

	if ferr == nil {
		if f.IsDir() {
			return value, SFOLDER
		} else if f.Mode().IsRegular() {
			return value, SFILE
		}
	}

	return value, SERROR
}

// ListConfig return a config list
func ListConfig() (*[]Config, error) {
	configs := []Config{}

	rows, err := db.Query("select d.id, d.category, d.identifier, d.created_at, s.id, s.category, s.identifier, s.created_at, s.doc_id from docs d inner join sources s on s.doc_id = d.id order by d.created_at")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		doc := Doc{}
		source := Source{}

		err := rows.Scan(&doc.ID, &doc.Category, &doc.Identifier, &doc.CreatedAt, &source.ID, &source.Category, &source.Identifier, &source.CreatedAt, &source.DocID)

		if err != nil {
			return nil, err
		}

		if len(configs) == 0 || configs[len(configs)-1].Doc.ID != source.DocID {
			configs = append(configs, Config{
				Doc: doc,
			})
		}

		configs[len(configs)-1].Sources = append(configs[len(configs)-1].Sources, source)
	}

	return &configs, nil
}

// CreateConfig populates everything needed to create a new config entry
func CreateConfig(doc *Doc, sources *[]Source) error {
	err := InsertDoc(doc)

	if err != nil {
		return err
	}

	for _, source := range *sources {
		switch source.Category {
		case SFILE:
			err = createFileSource(&source)
		case SFOLDER:
			err = createFolderSource(&source)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveConfigs delete one or several config group
func RemoveConfigs(configs *[]Config) error {
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
			return err
		}

		_, err = db.Exec(fmt.Sprintf("delete from docs where id in (%s)", strings.Join(docIds, ",")))

		if err != nil {
			return err
		}

		_, err = db.Exec(fmt.Sprintf("delete from items where source_id in (%s)", strings.Join(sourceIds, ",")))

		if err != nil {
			return err
		}
	}

	return nil
}
