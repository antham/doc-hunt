package model

import (
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

	content := []byte("<?php echo 'A source file';")
	err := ioutil.WriteFile("/tmp/doc-hunt/source1.php", content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	content = []byte("<?php echo 'Another source file';")
	err = ioutil.WriteFile("/tmp/doc-hunt/source2.php", content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}
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
