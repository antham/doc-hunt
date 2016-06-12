package file

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTables(t *testing.T) {
	id := uuid.NewV4().String()

	dbName = "/tmp/" + id

	createTables()

	_, err := db.Exec("select * from docs")

	assert.NoError(t, err, "Select * docs table must return no error")

	_, err = db.Exec("select * from sources")

	assert.NoError(t, err, "Select * sources table must return no error")
}
