package file

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func insertFakeConfig(docCat DocCategory, docIdentifier string, sourceIdentifiers map[string]SourceCategory) (Doc, []Source) {
	sources := []Source{}

	doc := NewDoc(docIdentifier, docCat)
	err := InsertDoc(doc)

	if err != nil {
		logrus.Fatal(err)
	}

	for sidentifier, scat := range sourceIdentifiers {
		source := NewSource(doc, sidentifier, scat)
		err := InsertSource(source)

		sources = append(sources, *source)

		if err != nil {
			logrus.Fatal(err)
		}
	}

	return *doc, sources
}

func TestNewSource(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	createSubTestDirectory("whatever")

	doc, sources := insertFakeConfig(DFILE, "test.txt", map[string]SourceCategory{"whatever": SFILEREG})

	assert.EqualValues(t, SFILEREG, sources[0].Category, "Must be a folder source")
	assert.Equal(t, doc.ID, sources[0].DocID, "Must be tied to document declared previously")
}

func TestInsertConfig(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	doc, sources := insertFakeConfig(DFILE, "doc_file_to_track.txt", map[string]SourceCategory{"source1.php": SFILEREG, "source2.php": SFILEREG})

	var identifier string

	err = db.QueryRow("select identifier from docs where id = ?", doc.ID).Scan(&identifier)

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
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	configs, err := ListConfig()

	assert.NoError(t, err, "Must return no errors")

	expected := []Config{}

	assert.Equal(t, &expected, configs, "Must return an empty config array")
}

func TestListConfigWithEntries(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	doc := NewDoc("doc_file_to_track.txt", DFILE)
	err = InsertDoc(doc)

	if err != nil {
		logrus.Fatal(err)
	}

	insertFakeConfig(DFILE, "doc_file_to_track.txt", map[string]SourceCategory{"source1.php": SFILEREG, "source2.php": SFILEREG})
	insertFakeConfig(DFILE, "doc_file_to_track.txt", map[string]SourceCategory{"source3.php": SFILEREG, "source4.php": SFILEREG, "source5.php": SFILEREG})

	configs, err := ListConfig()

	assert.NoError(t, err, "Must return no errors")

	assert.Len(t, *configs, 2, "Must have 2 configs")

	expectedResults := map[string]bool{
		"source1.php": true,
		"source2.php": true,
		"source3.php": true,
		"source4.php": true,
		"source5.php": true,
	}

	assert.Len(t, (*configs)[0].Sources, 2, "Must have 2 source files")
	assert.Len(t, (*configs)[1].Sources, 3, "Must have 3 source files")

	assert.Condition(t, func() bool {
		for i := 0; i < 2; i++ {
			if !expectedResults[(*configs)[0].Sources[i].Identifier] {
				return false
			}

			delete(expectedResults, (*configs)[0].Sources[i].Identifier)
		}

		for i := 0; i < 3; i++ {
			if !expectedResults[(*configs)[1].Sources[i].Identifier] {
				return false
			}

			delete(expectedResults, (*configs)[1].Sources[i].Identifier)
		}

		return true
	}, "Must return correct identifier")

	assert.Len(t, expectedResults, 0, "Must record all source files")
}

func TestRemoveConfigsWithNoResults(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	configs := []Config{}

	err = RemoveConfigs(&configs)
	if err != nil {
		logrus.Fatal(err)
	}

	result, err := ListConfig()

	assert.NoError(t, err, "Must return no errors")

	assert.Len(t, *result, 0, "Must have no config")
}

func TestRemoveConfigsWithOneEntry(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	insertFakeConfig(DFILE, "doc_file_to_track.txt", map[string]SourceCategory{"source1.php": SFILEREG, "source2.php": SFILEREG})

	configs, err := ListConfig()

	assert.NoError(t, err, "Must return no errors")

	err = RemoveConfigs(configs)
	if err != nil {
		logrus.Fatal(err)
	}

	result, err := ListConfig()

	assert.NoError(t, err, "Must return no errors")

	assert.Len(t, *result, 0, "Must have no config remaining")
}

func TestRemoveConfigsWithSeveralEntries(t *testing.T) {
	createMocks()
	deleteDatabase()

	err := Initialize()
	if err != nil {
		logrus.Fatal(err)
	}

	insertFakeConfig(DFILE, "doc_file_to_track.txt", map[string]SourceCategory{"source1.php": SFILEREG, "source2.php": SFILEREG})
	insertFakeConfig(DFILE, "doc_file_to_track.txt", map[string]SourceCategory{"source3.php": SFILEREG, "source4.php": SFILEREG})

	configs, err := ListConfig()
	assert.NoError(t, err, "Must return no errors")

	expected := (*configs)[1]

	c := append((*configs)[:1], (*configs)[2:]...)

	err = RemoveConfigs(&c)
	if err != nil {
		logrus.Fatal(err)
	}

	result, err := ListConfig()
	assert.NoError(t, err, "Must return no errors")

	assert.Len(t, *result, 1, "Must have no config remaining")
	assert.Equal(t, expected, (*result)[0], "Wrong configs deleted")
}

func TestNewItems(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	_, sources := insertFakeConfig(DFILE, "test.txt", map[string]SourceCategory{"whatever": SFILEREG})

	files := []string{"source1.php", "source2.php"}

	items, err := NewItems(&files, &sources[0])
	assert.NoError(t, err, "Must return no errors")

	assert.Equal(t, "source1.php", (*items)[0].Identifier, "Must return file source1.php")
	assert.Equal(t, "source2.php", (*items)[1].Identifier, "Must return file source2.php")
}
