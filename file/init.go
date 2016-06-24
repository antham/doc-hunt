package file

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

var db *sql.DB
var dbName = ".doc-hunt"

func init() {
	createTables()
}

func createTables() {
	var err error

	db, err = sql.Open("sqlite3", dbName)

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}

	createSourceTable()
	createDocTable()
}

func createDocTable() {
	query := `
create table docs(
id text primary key not null,
category int not null,
identifier text not null,
created_at timestamp not null);`

	_, err := db.Exec(query)

	if err != nil && err.(sqlite3.Error).Code != sqlite3.ErrError {
		ui.Error(err)

		util.ErrorExit()
	}
}

func createSourceTable() {
	query := `
create table sources(
id text primary key not null,
path text not null,
fingerprint text not null,
created_at timestamp not null,
updated_at timestamp not null,
doc_id text not null,
foreign key(doc_id) references doc_file(id));`

	_, err := db.Exec(query)

	if err != nil && err.(sqlite3.Error).Code != sqlite3.ErrError {
		ui.Error(err)

		util.ErrorExit()
	}
}
