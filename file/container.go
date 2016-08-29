package file

import (
	"database/sql"

	"github.com/antham/doc-hunt/util"
)

type container struct {
	db       *sql.DB
	dbName   string
	cfgRepo  *ConfigRepository
	srcRepo  *SourceRepository
	itemRepo *ItemRepository
	docRepo  *DocRepository
}

func newContainer(dbName string) (container, error) {
	db, err := sql.Open("sqlite3", util.GetAbsPath(dbName))

	if err != nil {
		return container{}, err
	}

	return container{
		dbName: dbName,
		db:     db,
	}, nil
}

func (c container) GetDatabase() *sql.DB {
	return c.db
}

func (c container) GetConfigRepository() *ConfigRepository {
	if c.cfgRepo == nil {
		repo := NewConfigRepository(c.GetDatabase(), c.GetSourceRepository(), c.GetItemRepository(), c.GetDocRepository())
		c.cfgRepo = &repo
	}

	return c.cfgRepo
}

func (c container) GetSourceRepository() *SourceRepository {
	if c.srcRepo == nil {
		repo := NewSourceRepository(c.GetDatabase())
		c.srcRepo = &repo
	}

	return c.srcRepo
}

func (c container) GetItemRepository() *ItemRepository {
	if c.itemRepo == nil {
		repo := NewItemRepository(c.GetDatabase())
		c.itemRepo = &repo
	}

	return c.itemRepo
}

func (c container) GetDocRepository() *DocRepository {
	if c.itemRepo == nil {
		repo := NewDocRepository(c.GetDatabase())
		c.docRepo = &repo
	}

	return c.docRepo
}
