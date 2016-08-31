package file

import (
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
	var err error

	res, err := Container.GetDatabase().Query("select id from version")

	defer func() {
		if e := res.Close(); e != nil {
			err = e
		}
	}()

	if err == nil && !res.Next() {
		_, err = Container.GetDatabase().Exec("insert into version values (?)", appVersion)

		return err
	}

	return err
}
