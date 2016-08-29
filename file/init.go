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

	err = createDocTable()

	if err != nil {
		return err
	}

	err = createSourceTable()

	if err != nil {
		return err
	}

	err = createItemTable()

	if err != nil {
		return err
	}

	err = createVersionTable()

	if err != nil {
		return err
	}

	return initVersion()
}

func createVersionTable() error {
	query := `
create table version(
id text primary key not null
);`
	_, err := Container.GetDatabase().Exec(query)

	if err != nil && err.(sqlite3.Error).Code != sqlite3.ErrError {
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

	_, err := Container.GetDatabase().Exec(query)

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

	_, err := Container.GetDatabase().Exec(query)

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
