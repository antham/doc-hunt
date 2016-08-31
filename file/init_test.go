package file

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateTables(t *testing.T) {
	removeTestDirectory()
	createTestDirectory()

	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	_, err = Container.GetDatabase().Exec("select * from docs")

	assert.NoError(t, err, "Select docs table must return no error")

	_, err = Container.GetDatabase().Exec("select * from sources")

	assert.NoError(t, err, "Select sources table must return no error")

	_, err = Container.GetDatabase().Exec("select * from items")

	assert.NoError(t, err, "Select items table must return no error")

	version, err := Container.GetSettingRepository().Get("version")

	assert.NoError(t, err, "Retrieves version must return no errors")
	assert.Equal(t, appVersion, version, "Must record version in settings")
}

func TestMoveVersionToSettings(t *testing.T) {
	removeTestDirectory()
	createTestDirectory()

	appVersion = "0.0.1"

	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	query := `
create table version(
id text primary key not null
);`
	_, err = Container.GetDatabase().Exec(query)

	if err != nil {
		logrus.Fatal(err)
	}

	_, err = Container.GetDatabase().Exec("insert into version values (?)", "4.5.6")

	if err != nil {
		logrus.Fatal(err)
	}

	_, err = Container.GetDatabase().Exec("delete from settings")

	if err != nil {
		logrus.Fatal(err)
	}

	err = initVersion()

	if err != nil {
		logrus.Fatal(err)
	}

	var version string

	err = Container.GetDatabase().QueryRow("select value from settings where name = (?)", "version").Scan(&version)

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "4.5.6", version, "Must return version previously recorded in version table")
}
