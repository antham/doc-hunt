package file

import (
	"fmt"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	createMocks()
	deleteDatabase()
	Initialize()
	createSubTestDirectory("whatever")

	doc := NewDoc("test.txt", DFILE)
	source := NewSource(doc, "whatever", SFOLDER)

	assert.EqualValues(t, SFOLDER, source.Category, "Must be a folder source")
	assert.Equal(t, doc.ID, source.DocID, "Must be tied to document declared previously")
}

func TestInsertConfig(t *testing.T) {
	createMocks()
	deleteDatabase()
	Initialize()

	doc := NewDoc("doc_file_to_track.txt", DFILE)
	InsertDoc(doc)

	sources := []*Source{NewSource(doc, "source1.php", SFILE), NewSource(doc, "source2.php", SFILE)}

	InsertSource(sources[0])
	InsertSource(sources[1])

	var identifier string

	err := db.QueryRow("select identifier from docs where id = ?", doc.ID).Scan(&identifier)

	assert.NoError(t, err, "Must return no errors")

	assert.Equal(t, doc.Identifier, identifier, "Must record a doc file row")

	results, err := db.Query("select identifier, category from sources where doc_id = ? order by identifier", doc.ID)

	assert.NoError(t, err, "Must return no errors")

	defer func() {
		if err := results.Close(); err != nil {
			logrus.Fatal(err)
		}
	}()

	var i int

	for results.Next() {
		var identifier string
		var category SourceCategory

		err := results.Scan(&identifier, &category)

		assert.NoError(t, err, "Must return no errors")

		assert.Equal(t, sources[i].Identifier, identifier, "Must return source identifier : "+sources[i].Identifier)
		assert.Equal(t, sources[i].Category, category, fmt.Sprintf("Must return fingerprint identifier : %d", sources[i].Category))

		i++
	}
}

func TestListConfigWithNoResults(t *testing.T) {
	createMocks()
	deleteDatabase()
	Initialize()

	configs, err := ListConfig()

	assert.NoError(t, err, "Must return no errors")

	expected := []Config{}

	assert.Equal(t, &expected, configs, "Must return an empty config array")
}

func TestListConfigWithEntries(t *testing.T) {
	createMocks()
	deleteDatabase()
	Initialize()

	doc := NewDoc("doc_file_to_track.txt", DFILE)
	InsertDoc(doc)

	sources := []*Source{NewSource(doc, "source1.php", SFILE), NewSource(doc, "source2.php", SFILE)}

	InsertSource(sources[0])
	InsertSource(sources[1])

	doc = NewDoc("doc_file_to_track.txt", DFILE)
	InsertDoc(doc)

	sources = []*Source{NewSource(doc, "source3.php", SFILE), NewSource(doc, "source4.php", SFILE), NewSource(doc, "source5.php", SFILE)}

	InsertSource(sources[0])
	InsertSource(sources[1])
	InsertSource(sources[2])

	configs, err := ListConfig()

	assert.NoError(t, err, "Must return no errors")

	assert.Len(t, *configs, 2, "Must have 2 configs")

	assert.Len(t, (*configs)[0].Sources, 2, "Must have 2 source files")
	assert.Equal(t, "source1.php", (*configs)[0].Sources[0].Identifier, "Must return correct identifier")
	assert.Equal(t, "source2.php", (*configs)[0].Sources[1].Identifier, "Must return correct identifier")

	assert.Len(t, (*configs)[1].Sources, 3, "Must have 3 source files")
	assert.Equal(t, "source3.php", (*configs)[1].Sources[0].Identifier, "Must return correct identifier")
	assert.Equal(t, "source4.php", (*configs)[1].Sources[1].Identifier, "Must return correct identifier")
	assert.Equal(t, "source5.php", (*configs)[1].Sources[2].Identifier, "Must return correct identifier")
}

func TestRemoveConfigsWithNoResults(t *testing.T) {
	createMocks()
	deleteDatabase()
	Initialize()

	configs := []Config{}

	RemoveConfigs(&configs)

	result, err := ListConfig()

	assert.NoError(t, err, "Must return no errors")

	assert.Len(t, *result, 0, "Must have no config")
}

func TestRemoveConfigsWithOneEntry(t *testing.T) {
	createMocks()
	deleteDatabase()
	Initialize()

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []*Source{NewSource(doc, "source1.php", SFILE), NewSource(doc, "source2.php", SFILE)}

	InsertSource(sources[0])
	InsertSource(sources[1])

	configs, err := ListConfig()

	assert.NoError(t, err, "Must return no errors")

	RemoveConfigs(configs)

	result, err := ListConfig()

	assert.NoError(t, err, "Must return no errors")

	assert.Len(t, *result, 0, "Must have no config remaining")
}

func TestRemoveConfigsWithSeveralEntries(t *testing.T) {
	createMocks()
	deleteDatabase()
	Initialize()

	doc := NewDoc("doc_file_to_track.txt", DFILE)
	InsertDoc(doc)

	InsertSource(NewSource(doc, "source1.php", SFILE))
	InsertSource(NewSource(doc, "source2.php", SFILE))

	doc = NewDoc("doc_file_to_track.txt", DFILE)
	InsertDoc(doc)

	InsertSource(NewSource(doc, "source3.php", SFILE))
	InsertSource(NewSource(doc, "source4.php", SFILE))

	doc = NewDoc("doc_file_to_track.txt", DFILE)
	InsertDoc(doc)

	InsertSource(NewSource(doc, "source1.php", SFILE))
	InsertSource(NewSource(doc, "source2.php", SFILE))

	configs, err := ListConfig()
	assert.NoError(t, err, "Must return no errors")

	expected := (*configs)[1]

	c := append((*configs)[:1], (*configs)[2:]...)

	RemoveConfigs(&c)

	result, err := ListConfig()
	assert.NoError(t, err, "Must return no errors")

	assert.Len(t, *result, 1, "Must have no config remaining")
	assert.Equal(t, expected, (*result)[0], "Wrong configs deleted")
}

func TestNewItems(t *testing.T) {
	createMocks()
	deleteDatabase()
	Initialize()

	doc := NewDoc("test.txt", DFILE)
	source := NewSource(doc, "whatever", SFOLDER)

	files := []string{"source1.php", "source2.php"}

	items, err := NewItems(&files, source)
	assert.NoError(t, err, "Must return no errors")

	assert.Equal(t, "source1.php", (*items)[0].Identifier, "Must return file source1.php")
	assert.Equal(t, "source2.php", (*items)[1].Identifier, "Must return file source2.php")
}
