package file

import (
	"database/sql"
)

// SettingRepository handle all actions availables on settings table
type SettingRepository struct {
	db *sql.DB
}

// NewSettingRepository return a new SettingRepository instance
func NewSettingRepository(db *sql.DB) SettingRepository {
	return SettingRepository{db: db}
}

// Get retrieve a setting entry
func (s SettingRepository) Get(name string) (string, error) {
	var value string

	err := s.db.QueryRow("select value from settings where name = (?)", name).Scan(&value)

	return value, err
}

// Create insert a new setting entry
func (s SettingRepository) Create(setting *Setting) error {
	_, err := s.db.Exec("insert into settings values (?,?)", setting.Name, setting.Value)

	return err
}

// Update set an already existing setting entry with a new value
func (s SettingRepository) Update(setting *Setting) error {
	_, err := s.db.Exec("update settings set value=(?) where name=(?)", setting.Value, setting.Name)

	return err
}
