package file

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTables(t *testing.T) {
	deleteDatabase()

	id := uuid.NewV4().String()

	Container.dbName = id

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

	_, err = Container.GetDatabase().Exec("select * from version")

	assert.NoError(t, err, "Select version table must return no error")
}
