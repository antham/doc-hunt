package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/util"
)

func retrieveItems(identifiers []string) map[string]*[]Item {
	items := map[string]*[]Item{}

	for _, identifier := range identifiers {
		var id string
		err := db.QueryRow("select id from sources where identifier = ?", identifier).Scan(&id)

		if err != nil {
			logrus.Warn(err)
		}

		source := Source{ID: id}
		items[identifier], err = getItems(&source)

		if err != nil {
			logrus.Warn(err)
		}
	}

	return items
}

func TestUpdateItemsFingerprint(t *testing.T) {
	createMocks()
	deleteDatabase()
	Initialize()

	err := CreateConfig("doc_file_to_track.txt", DFILE, []string{}, []string{"source1.php", "source2.php"})

	before := retrieveItems([]string{"source1.php", "source2.php"})

	if err != nil {
		logrus.Fatal(err)
	}

	err = ioutil.WriteFile(util.GetAbsPath("source1.php"), []byte("<?php echo 'Hello world !';"), 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	updateItemsFingeprint()

	after := retrieveItems([]string{"source1.php", "source2.php"})

	assert.True(t, (*after["source1.php"])[0].UpdatedAt.After((*before["source1.php"])[0].UpdatedAt), "Must changes updated date")

	assert.True(t, (*after["source2.php"])[0].UpdatedAt.After((*before["source1.php"])[0].UpdatedAt), "Must changes updated date")
}

func TestDeleteItems(t *testing.T) {
	createMocks()
	deleteDatabase()
	Initialize()

	err := CreateConfig("doc_file_to_track.txt", DFILE, []string{}, []string{"source1.php", "source2.php"})

	before := retrieveItems([]string{"source1.php", "source2.php"})

	if err != nil {
		logrus.Fatal(err)
	}

	err = os.Remove(util.GetAbsPath("source1.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	deleteItems(&[]string{"source1.php"})

	after := retrieveItems([]string{"source1.php", "source2.php"})

	assert.Len(t, (*before["source1.php"]), 1, "Must contains 1 element")
	assert.Len(t, (*before["source2.php"]), 1, "Must contains 1 element")

	assert.Len(t, (*after["source1.php"]), 0, "Must contains no element")
	assert.Len(t, (*after["source2.php"]), 1, "Must contains 1 element, only first item is removed")
}

func TestDeleteItemsWithOnlyOneItemRemaining(t *testing.T) {
	createMocks()
	createDocFiles()
	deleteDatabase()
	Initialize()

	err := CreateConfig("doc_file_to_track.txt", DFILE, []string{}, []string{"source1.php"})

	if err != nil {
		logrus.Fatal(err)
	}

	err = CreateConfig("doc_file_to_track_2.txt", DFILE, []string{}, []string{"source2.php"})

	if err != nil {
		logrus.Fatal(err)
	}

	before := retrieveItems([]string{"source1.php", "source2.php"})

	err = os.Remove(util.GetAbsPath("source1.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	deleteItems(&[]string{"source1.php"})

	sourceRows, err := db.Query("select s.id from sources s where identifier = ?", "source1.php")

	if err != nil {
		logrus.Fatal(err)
	}

	docRows, err := db.Query("select d.id from docs d where identifier = ?", "doc_file_to_track.txt")

	if err != nil {
		logrus.Fatal(err)
	}

	after := retrieveItems([]string{"source1.php", "source2.php"})

	assert.Len(t, (*before["source1.php"]), 1, "Must contains 1 element")
	assert.Len(t, (*before["source2.php"]), 1, "Must contains 1 element")

	assert.Len(t, (*after["source1.php"]), 0, "Must contains no element")
	assert.Len(t, (*after["source2.php"]), 1, "Must contains 1 element, only first item is removed")

	assert.False(t, sourceRows.Next(), "Must have deleted source id")
	assert.False(t, docRows.Next(), "Must have deleted doc id")
}

func TestUpdate(t *testing.T) {
	createMocks()
	createDocFiles()
	deleteDatabase()
	Initialize()
	createSubTestDirectory("test1")
	createSourceFilesInPath("test1")

	err := CreateConfig("doc_file_to_track.txt", DFILE, []string{"test1"}, []string{})

	if err != nil {
		logrus.Fatal(err)
	}

	err = os.Remove(util.GetAbsPath("test1/source1.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	content := []byte("whatever")
	err = ioutil.WriteFile(util.GetAbsPath("test1/source2.php"), content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	createSourceFile([]byte("test"), "test1/source20.php")

	err = Update()

	assert.NoError(t, err, "Must produces no errors")

	items := retrieveItems([]string{"test1"})

	assert.Len(t, (*items["test1"]), 10, "Must return 10 items")

	values := map[string]Item{}

	for _, i := range *items["test1"] {
		values[i.Identifier] = i
	}

	_, exists := values["test1/source1.php"]

	assert.Equal(t, values["test1/source20.php"].Fingerprint, "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", "Must a new files")
	assert.Equal(t, values["test1/source2.php"].Fingerprint, "d869db7fe62fb07c25a0403ecaea55031744b5fb", "Must update fingerprint of changed files")
	assert.False(t, exists, "Must remove deleted file")
}
