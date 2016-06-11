package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func createTestDirectory() {
	os.RemoveAll("/tmp/doc-hunt")

	err := os.Mkdir("/tmp/doc-hunt", 0777)

	if err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func createADocFile() {
	createTestDirectory()

	content := []byte("This is a doc file")
	err := ioutil.WriteFile("/tmp/doc-hunt/doc_file_to_track.txt", content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}
}

func createSourceFiles() {
	createTestDirectory()

	for i := 1; i <= 10; i++ {
		content := []byte("<?php echo 'A source file';")
		err := ioutil.WriteFile(fmt.Sprintf("/tmp/doc-hunt/source%d.php", i), content, 0644)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func deleteDatabase() {
	os.Remove(dbName)
}

func TestInsertConfig(t *testing.T) {
	createADocFile()
	createSourceFiles()

	doc := NewDoc("/tmp/doc-hunt/doc_file_to_track.txt")
	sources := NewSources(doc, []string{"/tmp/doc-hunt/source1.php", "/tmp/doc-hunt/source2.php"})

	InsertConfig(doc, sources)

	var path string

	err := db.QueryRow("select path from doc_file where id = ?", doc.ID).Scan(&path)

	assert.NoError(t, err, "Must return no errors")

	assert.Equal(t, doc.Path, path, "Must record a doc file row")

	results, err := db.Query("select path, fingerprint from source_file where doc_file_id = ? order by path", doc.ID)

	assert.NoError(t, err, "Must return no errors")

	defer results.Close()

	var i int

	for results.Next() {
		var path string
		var fingerprint string

		err := results.Scan(&path, &fingerprint)

		assert.NoError(t, err, "Must return no errors")

		assert.Equal(t, (*sources)[i].Path, path, "Must return source path : "+(*sources)[i].Path)
		assert.Equal(t, (*sources)[i].Fingerprint, fingerprint, "Must return fingerprint path : "+(*sources)[i].Fingerprint)

		i++
	}
}

func TestListConfigWithNoResults(t *testing.T) {
	deleteDatabase()
	createTables()

	configs := ListConfig()

	expected := []Config{}

	assert.Equal(t, &expected, configs, "Must return an empty config array")
}

func TestListConfigWithEntries(t *testing.T) {
	deleteDatabase()
	createTables()

	doc := NewDoc("/tmp/doc-hunt/doc_file_to_track.txt")
	sources := NewSources(doc, []string{"/tmp/doc-hunt/source1.php", "/tmp/doc-hunt/source2.php"})

	InsertConfig(doc, sources)

	doc = NewDoc("/tmp/doc-hunt/doc_file_to_track.txt")
	sources = NewSources(doc, []string{"/tmp/doc-hunt/source3.php", "/tmp/doc-hunt/source4.php", "/tmp/doc-hunt/source5.php"})

	InsertConfig(doc, sources)

	configs := ListConfig()

	assert.Len(t, *configs, 2, "Must have 2 configs")

	assert.Len(t, (*configs)[0].SourceFiles, 2, "Must have 2 source files")
	assert.Equal(t, "/tmp/doc-hunt/source1.php", (*configs)[0].SourceFiles[0].Path, "Must return correct path")
	assert.Equal(t, "/tmp/doc-hunt/source2.php", (*configs)[0].SourceFiles[1].Path, "Must return correct path")

	assert.Len(t, (*configs)[1].SourceFiles, 3, "Must have 3 source files")
	assert.Equal(t, "/tmp/doc-hunt/source3.php", (*configs)[1].SourceFiles[0].Path, "Must return correct path")
	assert.Equal(t, "/tmp/doc-hunt/source4.php", (*configs)[1].SourceFiles[1].Path, "Must return correct path")
	assert.Equal(t, "/tmp/doc-hunt/source5.php", (*configs)[1].SourceFiles[2].Path, "Must return correct path")
}

func TestRemoveConfigsWithNoResults(t *testing.T) {
	deleteDatabase()
	createTables()

	configs := []Config{}

	RemoveConfigs(&configs)

	result := ListConfig()

	assert.Len(t, *result, 0, "Must have no config")
}

func TestRemoveConfigsWithOneEntry(t *testing.T) {
	deleteDatabase()
	createTables()

	doc := NewDoc("/tmp/doc-hunt/doc_file_to_track.txt")
	sources := NewSources(doc, []string{"/tmp/doc-hunt/source1.php", "/tmp/doc-hunt/source2.php"})

	InsertConfig(doc, sources)

	configs := ListConfig()

	RemoveConfigs(configs)

	result := ListConfig()

	assert.Len(t, *result, 0, "Must have no config remaining")
}

func TestRemoveConfigsWithSeveralEntries(t *testing.T) {
	deleteDatabase()
	createTables()

	doc := NewDoc("/tmp/doc-hunt/doc_file_to_track.txt")
	sources := NewSources(doc, []string{"/tmp/doc-hunt/source1.php", "/tmp/doc-hunt/source2.php"})

	InsertConfig(doc, sources)

	doc = NewDoc("/tmp/doc-hunt/doc_file_to_track.txt")
	sources = NewSources(doc, []string{"/tmp/doc-hunt/source3.php", "/tmp/doc-hunt/source4.php"})

	InsertConfig(doc, sources)

	doc = NewDoc("/tmp/doc-hunt/doc_file_to_track.txt")
	sources = NewSources(doc, []string{"/tmp/doc-hunt/source3.php", "/tmp/doc-hunt/source5.php"})

	InsertConfig(doc, sources)

	configs := ListConfig()
	expected := (*configs)[1]

	c := append((*configs)[:1], (*configs)[2:]...)

	RemoveConfigs(&c)

	result := ListConfig()

	assert.Len(t, *result, 1, "Must have no config remaining")
	assert.Equal(t, expected, (*result)[0], "Wrong configs deleted")
}
