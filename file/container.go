package file

import (
	"database/sql"
	"os"

	"github.com/antham/doc-hunt/util"
)

type container struct {
	db          *sql.DB
	dbName      string
	cfgRepo     *ConfigRepository
	srcRepo     *SourceRepository
	itemRepo    *ItemRepository
	docRepo     *DocRepository
	settingRepo *SettingRepository
	manager     *Manager
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

func (c container) GetSettingRepository() *SettingRepository {
	if c.itemRepo == nil {
		repo := NewSettingRepository(c.GetDatabase())
		c.settingRepo = &repo
	}

	return c.settingRepo
}

func (c container) GetManager() *Manager {
	if c.manager == nil {
		manager := NewManager(c.GetConfigRepository(), c.GetItemRepository(), util.GetAbsPath, util.ExtractFilesMatchingReg, os.Stat)
		c.manager = &manager
	}

	return c.manager
}
