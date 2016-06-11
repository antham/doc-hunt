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

	_, err := db.Exec("select * from doc_file")

	assert.NoError(t, err, "Select * doc_file table must return no error")

	_, err = db.Exec("select * from source_file")

	assert.NoError(t, err, "Select * source_file table must return no error")
}
