package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/util"
)

func TestFetchStatusWithSourceFilesUntouched(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"source1.php", "source2.php"})

	InsertConfig(doc, sources)

	results := FetchStatus()

	assert.Len(t, *results, 1, "One configuration must be returned")
	assert.Equal(t, "doc_file_to_track.txt", (*results)[0].Doc.Identifier, "Wrong doc identifier returned")
	assert.Len(t, (*results)[0].Sources, 2, "Two file sources must be found")
	assert.Equal(t, "source1.php", (*results)[0].Sources[0].Path, "Wrong source path returned")
	assert.Equal(t, "source2.php", (*results)[0].Sources[1].Path, "Wrong source path returned")
	assert.Equal(t, []string{"source1.php", "source2.php"}, (*results)[0].Status[Untouched], "Files must be marked as untouched")

}

func TestFetchStatusWithUpdatedSource(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"source1.php", "source2.php"})

	InsertConfig(doc, sources)

	content := []byte("whatever")
	err := ioutil.WriteFile(util.GetAbsPath("source1.php"), content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	results := FetchStatus()

	assert.Len(t, *results, 1, "One configuration must be returned")
	assert.Equal(t, (*results)[0].Status[Updated], []string{"source1.php"}, "File must be marked as Deleted")
	assert.Equal(t, (*results)[0].Status[Untouched], []string{"source2.php"}, "File must be marked as untouched")
}

func TestFetchStatusWithDeletedSource(t *testing.T) {
	createMocks()
	deleteDatabase()
	createTables()

	doc := NewDoc("doc_file_to_track.txt", FILE)
	sources := NewSources(doc, []string{"source1.php", "source2.php"})

	InsertConfig(doc, sources)

	err := os.Remove(util.GetAbsPath("source1.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	results := FetchStatus()

	assert.Len(t, *results, 1, "One configuration must be returned")
	assert.Equal(t, (*results)[0].Status[Deleted], []string{"source1.php"}, "File must be marked as Deleted")
	assert.Equal(t, (*results)[0].Status[Untouched], []string{"source2.php"}, "File must be marked as untouched")
}
