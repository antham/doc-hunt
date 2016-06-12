package file

import (
	"database/sql"

	"github.com/Sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
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
		logrus.Fatal(err)
	}

	createSourceFileTable()
	createDocFileTable()
}

func createDocFileTable() {
	query := `
create table docs(
id text primary key not null,
category int not null,
identifier text not null,
created_at timestamp not null);`

	db.Exec(query)
}

func createSourceFileTable() {
	query := `
create table sources(
id text primary key not null,
path text not null,
fingerprint text not null,
created_at timestamp not null,
updated_at timestamp not null,
doc_id text not null,
foreign key(doc_id) references doc_file(id));`
	db.Exec(query)
}
