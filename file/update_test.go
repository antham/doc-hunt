package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/util"
)

func TestUpdateItemsFingerprint(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("doc_file_to_track.txt", DFILE)
	InsertDoc(doc)

	source1 := NewSource(doc, "source1.php", SFILE)
	source2 := NewSource(doc, "source2.php", SFILE)

	InsertSource(source1)
	InsertSource(source2)

	InsertItems(NewItems(&[]string{"source1.php"}, source1))
	InsertItems(NewItems(&[]string{"source2.php"}, source2))

	items1Before := getItems(source1)
	items2Before := getItems(source2)

	err := ioutil.WriteFile(util.GetAbsPath("source1.php"), []byte("<?php echo 'Hello world !';"), 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	UpdateItemsFingeprint()

	items1After := getItems(source1)
	items2After := getItems(source2)

	assert.Equal(t, (*items1Before)[0].ID, (*items1After)[0].ID, "Must return same id")
	assert.True(t, (*items1After)[0].UpdatedAt.After((*items1Before)[0].UpdatedAt), "Must changes updated date")

	assert.Equal(t, (*items2Before)[0].ID, (*items2After)[0].ID, "Must return same id")
	assert.True(t, (*items2After)[0].UpdatedAt.After((*items2Before)[0].UpdatedAt), "Must changes updated date")
}

func TestDeleteItems(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("doc_file_to_track.txt", DFILE)
	InsertDoc(doc)

	source1 := NewSource(doc, "source1.php", SFILE)
	source2 := NewSource(doc, "source2.php", SFILE)

	InsertItems(NewItems(&[]string{"source1.php"}, source1))
	InsertItems(NewItems(&[]string{"source2.php"}, source2))

	items1Before := getItems(source1)
	items2Before := getItems(source2)

	err := os.Remove(util.GetAbsPath("source1.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	DeleteItems(&[]string{"source1.php"})

	items1After := getItems(source1)
	items2After := getItems(source2)

	assert.Len(t, (*items1Before), 1, "Must contains 1 element")
	assert.Len(t, (*items2Before), 1, "Must contains 1 element")

	assert.Len(t, (*items1After), 0, "Must contains no element")
	assert.Len(t, (*items2After), 1, "Must contains 1 element, only first item is removed")
}

func TestDeleteItemsWithOnlyOneItemRemaining(t *testing.T) {
	createMocks()
	createDocFiles()
	deleteDatabase()
	createTables()

	doc1 := NewDoc("doc_file_to_track.txt", DFILE)
	InsertDoc(doc1)

	source1 := NewSource(doc1, "source1.php", SFILE)
	InsertSource(source1)

	doc2 := NewDoc("doc_file_to_track_2.txt", DFILE)
	InsertDoc(doc2)

	source2 := NewSource(doc2, "source2.php", SFILE)
	InsertSource(source2)

	InsertItems(NewItems(&[]string{"source1.php"}, source1))
	InsertItems(NewItems(&[]string{"source2.php"}, source2))

	items1Before := getItems(source1)
	items2Before := getItems(source2)

	err := os.Remove(util.GetAbsPath("source1.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	DeleteItems(&[]string{"source1.php"})

	sourceRows, err := db.Query("select s.id from sources s where id = ?", source1.ID)

	if err != nil {
		logrus.Fatal(err)
	}

	docRows, err := db.Query("select d.id from docs d where id = ?", doc1.ID)

	if err != nil {
		logrus.Fatal(err)
	}

	items1After := getItems(source1)
	items2After := getItems(source2)

	assert.Len(t, (*items1Before), 1, "Must contains 1 element")
	assert.Len(t, (*items2Before), 1, "Must contains 1 element")

	assert.Len(t, (*items1After), 0, "Must contains no element")
	assert.Len(t, (*items2After), 1, "Must contains 1 element, only first item is removed")

	assert.False(t, sourceRows.Next(), "Must have deleted source id")
	assert.False(t, docRows.Next(), "Must have deleted doc id")
}
