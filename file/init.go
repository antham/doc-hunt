package file

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"

	"github.com/antham/doc-hunt/util"
)

var db *sql.DB
var dbName = ".doc-hunt"

// Initialize initialize project
func Initialize() error {
	var err error

	db, err = sql.Open("sqlite3", util.GetAbsPath(dbName))

	if err != nil {
		return err
	}

	err = createSourceTable()

	if err != nil {
		return err
	}

	err = createDocTable()

	if err != nil {
		return err
	}

	err = createItemTable()

	if err != nil {
		return err
	}

	return nil
}

func createDocTable() error {
	query := `
create table docs(
id text primary key not null,
category int not null,
identifier text not null,
created_at timestamp not null);`

	_, err := db.Exec(query)

	if err != nil && err.(sqlite3.Error).Code != sqlite3.ErrError {
		return err
	}

	return nil
}

func createSourceTable() error {
	query := `
create table sources(
id text primary key not null,
identifier text not null,
category int not null,
created_at timestamp not null,
doc_id text not null,
foreign key(doc_id) references docs(id));`

	_, err := db.Exec(query)

	if err != nil && err.(sqlite3.Error).Code != sqlite3.ErrError {
		return err
	}

	return nil
}

func createItemTable() error {
	query := `
create table items(
id text primary key not null,
identifier text not null,
fingerprint text not null,
created_at timestamp not null,
updated_at timestamp not null,
source_id text not null,
foreign key(source_id) references sources(id));`

	_, err := db.Exec(query)

	if err != nil && err.(sqlite3.Error).Code != sqlite3.ErrError {
		return err
	}

	return nil
}
