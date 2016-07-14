package file

import (
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/util"
)

func TestBuildStatusWithFilesUntouched(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "source1.php", SFILE),
		*NewSource(doc, "source2.php", SFILE),
	}

	err = CreateConfig(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	results, err := BuildStatus()

	assert.NoError(t, err, "Must returns no errors")

	assert.Len(t, *results, 2, "Two results must be returned")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[0].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "source1.php", (*results)[0].Source.Identifier, "Wrong source identifier returned")
	assert.Len(t, (*results)[0].Status[INONE], 1, "One file item must be found")
	assert.Equal(t, []string{"source1.php"}, (*results)[0].Status[INONE], "One file item must be found")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[1].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "source2.php", (*results)[1].Source.Identifier, "Wrong source identifier returned")
	assert.Len(t, (*results)[1].Status[INONE], 1, "One file item must be found")
	assert.Equal(t, []string{"source2.php"}, (*results)[1].Status[INONE], "One file item must be found")
}

func TestBuildStatusWithUpdatedFile(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "source1.php", SFILE),
		*NewSource(doc, "source2.php", SFILE),
	}

	err = CreateConfig(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	content := []byte("whatever")
	err = ioutil.WriteFile(util.GetAbsPath("source1.php"), content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	results, err := BuildStatus()

	assert.NoError(t, err, "Must returns no errors")

	assert.Len(t, *results, 2, "Two results must be returned")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[0].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "source1.php", (*results)[0].Source.Identifier, "Wrong source identifier returned")
	assert.Len(t, (*results)[0].Status[IUPDATED], 1, "One file item must be found")
	assert.Equal(t, []string{"source1.php"}, (*results)[0].Status[IUPDATED], "One file item must be found")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[1].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "source2.php", (*results)[1].Source.Identifier, "Wrong source identifier returned")
	assert.Len(t, (*results)[1].Status[INONE], 1, "One file item must be found")
	assert.Equal(t, []string{"source2.php"}, (*results)[1].Status[INONE], "One file item must be found")
}

func TestBuildStatusWithDeletedFile(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "source1.php", SFILE),
		*NewSource(doc, "source2.php", SFILE),
	}

	err = CreateConfig(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	err = os.Remove(util.GetAbsPath("source1.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	results, err := BuildStatus()

	assert.NoError(t, err, "Must returns no errors")

	assert.Len(t, *results, 2, "Two results must be returned")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[0].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "source1.php", (*results)[0].Source.Identifier, "Wrong source identifier returned")
	assert.Len(t, (*results)[0].Status[IDELETED], 1, "One file item must be found")
	assert.Equal(t, []string{"source1.php"}, (*results)[0].Status[IDELETED], "One file item must be found")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[1].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "source2.php", (*results)[1].Source.Identifier, "Wrong source identifier returned")
	assert.Len(t, (*results)[1].Status[INONE], 1, "One file item must be found")
	assert.Equal(t, []string{"source2.php"}, (*results)[1].Status[INONE], "One file item must be found")
}

func TestBuildStatusWithFolderSource(t *testing.T) {
	createMocks()
	deleteDatabase()
	createSubTestDirectory("test1")
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	createSourceFile([]byte("test"), "test1/source1.php")
	createSourceFile([]byte("test"), "test1/source2.php")
	createSourceFile([]byte("test"), "test1/source3.php")

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "test1", SFOLDER),
	}

	err = CreateConfig(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	results, err := BuildStatus()

	assert.NoError(t, err, "Must returns no errors")

	sort.Strings((*results)[0].Status[INONE])

	assert.Len(t, *results, 1, "One result must be returned")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[0].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "test1", (*results)[0].Source.Identifier, "Wrong source identifier returned")
	assert.Len(t, (*results)[0].Status[INONE], 3, "Three file items must be found")
	assert.Equal(t, []string{"test1/source1.php", "test1/source2.php", "test1/source3.php"}, (*results)[0].Status[INONE], "Three file items must be found")
}

func TestBuildStatusWithFolderSourceAndAddedFile(t *testing.T) {
	createMocks()
	deleteDatabase()
	createSubTestDirectory("test1")
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	createSourceFile([]byte("test"), "test1/source1.php")
	createSourceFile([]byte("test"), "test1/source2.php")
	createSourceFile([]byte("test"), "test1/source3.php")

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "test1", SFOLDER),
	}

	err = CreateConfig(doc, &sources)

	createSourceFile([]byte("test"), "test1/source4.php")

	results, err := BuildStatus()

	assert.NoError(t, err, "Must returns no errors")

	sort.Strings((*results)[0].Status[INONE])

	assert.Len(t, *results, 1, "One result must be returned")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[0].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "test1", (*results)[0].Source.Identifier, "Wrong source identifier returned")

	assert.Len(t, (*results)[0].Status[INONE], 3, "Three file items must be found")
	assert.Equal(t, []string{"test1/source1.php", "test1/source2.php", "test1/source3.php"}, (*results)[0].Status[INONE], "Three files items must be found")

	assert.Len(t, (*results)[0].Status[IADDED], 1, "One file items must be found")
	assert.Equal(t, []string{"test1/source4.php"}, (*results)[0].Status[IADDED], "One file item must be found")
}

func TestBuildStatusWithFolderDeleted(t *testing.T) {
	createMocks()
	deleteDatabase()
	createSubTestDirectory("test1")
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	createSourceFile([]byte("test"), "test1/source1.php")

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "test1", SFOLDER),
	}

	err = CreateConfig(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	err = os.RemoveAll(util.GetAbsPath("test1"))

	if err != nil {
		logrus.Fatal(err)
	}

	results, err := BuildStatus()

	assert.NoError(t, err, "Must returns no errors")

	sort.Strings((*results)[0].Status[IDELETED])

	assert.Len(t, *results, 1, "One result must be returned")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[0].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "test1", (*results)[0].Source.Identifier, "Wrong source identifier returned")

	assert.Len(t, (*results)[0].Status[IDELETED], 1, "One file item must be found")
	assert.Equal(t, []string{"test1/source1.php"}, (*results)[0].Status[IDELETED], "One file item must be found")
}
