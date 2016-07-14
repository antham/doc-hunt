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
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	createDocFiles()
	createSourceFilesInPath("test1")

	err = CreateConfig("doc_file_to_track.txt", DFILE, []string{"test1"}, []string{"test1/source5.php"})

	if err != nil {
		logrus.Fatal(err)
	}

	err = os.Remove(util.GetAbsPath("test1/source2.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	content := []byte("whatever")
	err = ioutil.WriteFile(util.GetAbsPath("test1/source3.php"), content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	err = CreateConfig("doc_file_to_track_2.txt", DFILE, []string{}, []string{"test1/source6.php"})

	if err != nil {
		logrus.Fatal(err)
	}

	createSourceFile([]byte("test"), "test1/source11.php")

	itemStatus, changesOccured, err := GetItemStatus()

	assert.NoError(t, err, "Must return item status without errors")

	results := map[string]map[ItemStatus]map[string]bool{}

	for doc, status := range *itemStatus {
		if results[doc.Identifier] == nil {
			results[doc.Identifier] = map[ItemStatus]map[string]bool{
				IADDED:   status[IADDED],
				IUPDATED: status[IUPDATED],
				IDELETED: status[IDELETED],
				IFAILED:  status[IFAILED],
				INONE:    status[INONE],
			}
		}
	}

	assert.True(t, changesOccured, "Must indicate that some items changed")

	assert.Len(t, results["doc_file_to_track.txt"][IADDED], 1, "Must contains 1 element added")
	assert.Len(t, results["doc_file_to_track.txt"][IDELETED], 1, "Must contains 1 element deleted")
	assert.Len(t, results["doc_file_to_track.txt"][IUPDATED], 1, "Must contains 1 element updated")
	assert.Len(t, results["doc_file_to_track.txt"][INONE], 8, "Must contains 1 element untouched")

	assert.Len(t, results["doc_file_to_track_2.txt"][IADDED], 0, "Must contains no element added")
	assert.Len(t, results["doc_file_to_track_2.txt"][IDELETED], 0, "Must contains no element deleted")
	assert.Len(t, results["doc_file_to_track_2.txt"][IUPDATED], 0, "Must contains no element updated")
	assert.Len(t, results["doc_file_to_track_2.txt"][INONE], 1, "Must contains 1 element untouched")
}

func TestGetItemStatusWithNoChanges(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	createDocFiles()

	err = CreateConfig("doc_file_to_track.txt", DFILE, []string{}, []string{"source1.php"})

	if err != nil {
		logrus.Fatal(err)
	}

	itemStatus, changesOccured, err := GetItemStatus()

	assert.NoError(t, err, "Must return item status without errors")

	results := map[string]map[ItemStatus]map[string]bool{}

	for doc, status := range *itemStatus {
		if results[doc.Identifier] == nil {
			results[doc.Identifier] = map[ItemStatus]map[string]bool{
				IADDED:   status[IADDED],
				IUPDATED: status[IUPDATED],
				IDELETED: status[IDELETED],
				IFAILED:  status[IFAILED],
				INONE:    status[INONE],
			}
		}
	}

	assert.False(t, changesOccured, "Must indicate that no items changed")

	assert.Len(t, results["doc_file_to_track.txt"][IADDED], 0, "Must contains no element added")
	assert.Len(t, results["doc_file_to_track.txt"][IDELETED], 0, "Must contains no element deleted")
	assert.Len(t, results["doc_file_to_track.txt"][IUPDATED], 0, "Must contains no element updated")
	assert.Len(t, results["doc_file_to_track.txt"][INONE], 1, "Must contains 1 element untouched")
}
