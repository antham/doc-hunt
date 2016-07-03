package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/antham/doc-hunt/util"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetItemStatus(t *testing.T) {
	createMocks()
	deleteDatabase()
	createSubTestDirectory("test1")
	createSubTestDirectory("test2")
	createDocFiles()
	createTables()
	createDocFiles()

	doc1 := NewDoc("doc_file_to_track.txt", DFILE)
	InsertDoc(doc1)

	source1 := NewSource(doc1, "source5.php", SFOLDER)
	InsertSource(source1)

	InsertItems(NewItems(&[]string{"source5.php"}, source1))

	source2 := NewSource(doc1, "test1", SFOLDER)
	InsertSource(source2)

	InsertItems(NewItems(&[]string{"source1.php"}, source2))
	InsertItems(NewItems(&[]string{"source2.php"}, source2))
	InsertItems(NewItems(&[]string{"source3.php"}, source2))

	createSourceFile([]byte("test"), "test1/source4.php")

	err := os.Remove(util.GetAbsPath("source2.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	content := []byte("whatever")
	err = ioutil.WriteFile(util.GetAbsPath("source3.php"), content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	doc2 := NewDoc("doc_file_to_track_2.txt", DFILE)
	InsertDoc(doc2)

	source3 := NewSource(doc2, "source6.php", SFILE)
	InsertSource(source3)

	InsertItems(NewItems(&[]string{"source6.php"}, source3))

	itemStatus := GetItemStatus()

	results := map[string]map[ItemStatus]map[string]bool{}

	for doc, status := range *itemStatus {
		if results[doc.ID] == nil {
			results[doc.ID] = map[ItemStatus]map[string]bool{
				IADDED:   status[IADDED],
				IUPDATED: status[IUPDATED],
				IDELETED: status[IDELETED],
				IFAILED:  status[IFAILED],
				INONE:    status[INONE],
			}
		}
	}

	assert.Len(t, results[doc1.ID][IADDED], 1, "Must contains 1 element added")
	assert.Len(t, results[doc1.ID][IDELETED], 1, "Must contains 1 element deleted")
	assert.Len(t, results[doc1.ID][IUPDATED], 1, "Must contains 1 element updated")
	assert.Len(t, results[doc1.ID][INONE], 2, "Must contains 1 element untouched")

	assert.Len(t, results[doc2.ID][IADDED], 0, "Must contains no element added")
	assert.Len(t, results[doc2.ID][IDELETED], 0, "Must contains no element deleted")
	assert.Len(t, results[doc2.ID][IUPDATED], 0, "Must contains no element updated")
	assert.Len(t, results[doc2.ID][INONE], 1, "Must contains 1 element untouched")
}
