package file

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/antham/doc-hunt/util"
)

// ConfigRepository handle all actions availables on docs table
type ConfigRepository struct {
	db       *sql.DB
	srcRepo  *SourceRepository
	itemRepo *ItemRepository
	docRepo  *DocRepository
}

// NewConfigRepository return a new ConfigRepository instance
func NewConfigRepository(db *sql.DB, srcRepo *SourceRepository, itemRepo *ItemRepository, docRepo *DocRepository) ConfigRepository {
	return ConfigRepository{db: db, srcRepo: srcRepo, itemRepo: itemRepo, docRepo: docRepo}
}

func (c ConfigRepository) createFileRegSource(source *Source) error {
	err := c.srcRepo.Create(source)

	if err != nil {
		return err
	}

	files, err := util.ExtractFilesMatchingReg(source.Identifier)

	if err != nil {
		return err
	}

	return c.itemRepo.CreateFromIdentifiersAndSource(files, source)
}

// CreateFromDocAndSources insert a new config entry
func (c ConfigRepository) CreateFromDocAndSources(doc *Doc, sources *[]Source) error {
	r := c.docRepo
	err := r.Create(doc)

	if err != nil {
		return err
	}

	for _, source := range *sources {
		switch source.Category {
		case SFILEREG:
			err = c.createFileRegSource(&source)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// List return all availables config
func (c ConfigRepository) List() (*[]Config, error) {
	configs := []Config{}

	rows, err := c.db.Query("select d.id, d.category, d.identifier, d.created_at, s.id, s.category, s.identifier, s.created_at, s.doc_id from docs d inner join sources s on s.doc_id = d.id order by d.created_at")

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

// Delete remove one or several config group
func (c ConfigRepository) Delete(configs *[]Config) error {
	sourceIds := []string{}
	docIds := []string{}

	for _, config := range *configs {
		for _, source := range config.Sources {
			sourceIds = append(sourceIds, fmt.Sprintf(`"%s"`, source.ID))
		}
		docIds = append(docIds, fmt.Sprintf(`"%s"`, config.Doc.ID))
	}

	if len(sourceIds) > 0 {
		_, err := c.db.Exec(fmt.Sprintf("delete from sources where id in (%s)", strings.Join(sourceIds, ",")))

		if err != nil {
			return err
		}

		_, err = c.db.Exec(fmt.Sprintf("delete from docs where id in (%s)", strings.Join(docIds, ",")))

		if err != nil {
			return err
		}

		_, err = c.db.Exec(fmt.Sprintf("delete from items where source_id in (%s)", strings.Join(sourceIds, ",")))

		if err != nil {
			return err
		}
	}

	return nil
}
