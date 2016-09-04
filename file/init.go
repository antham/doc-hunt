package file

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"
)

// Container is a dependency injection container in charge of
// creating and keeping reference to objects
var Container container

var dbName = ".doc-hunt"

// Initialize initialize project
func Initialize() error {
	var err error

	Container, err = newContainer(dbName)

	if err != nil {
		return err
	}

	queries := []string{
		`create table docs(
id text primary key not null,
category int not null,
identifier text not null,
created_at timestamp not null);`,
		`create table sources(
id text primary key not null,
identifier text not null,
category int not null,
created_at timestamp not null,
doc_id text not null,
foreign key(doc_id) references docs(id));`,
		`create table items(
			id text primary key not null,
			identifier text not null,
			fingerprint text not null,
			created_at timestamp not null,
			updated_at timestamp not null,
			source_id text not null,
			foreign key(source_id) references sources(id));`,
		`create table settings(
			name text primary key,
			value text);`,
		`alter table items add vcs_ref text;`,
	}

	for _, query := range queries {
		err := runQuery(query)

		if err != nil {
			return err
		}
	}

	return initVersion()
}

func runQuery(query string) error {
	_, err := Container.GetDatabase().Exec(query)

	if err != nil && err.(sqlite3.Error).Code != sqlite3.ErrError {
		return err
	}

	return nil
}

func initVersion() error {
	var version, sversion, vversion, table string
	var err error

	// @deprecated : we store version in settings
	_ = Container.GetDatabase().QueryRow("select name from sqlite_master where type=(?) and name=(?)", "table", "version").Scan(&table)

	if table != "" {
		// @deprecated : we store version in settings
		if err = Container.GetDatabase().QueryRow("select id from version").Scan(&vversion); err != nil {
			return err
		}

		if vversion != "" {
			version = vversion
		}
	} else {
		if sversion, err = Container.GetSettingRepository().Get("version"); err != nil && err != sql.ErrNoRows {
			return err
		}

		if sversion == "" {
			version = appVersion
		}
	}

	if version != "" {
		return Container.GetSettingRepository().Create(NewSetting("version", version))
	}

	return nil
}
