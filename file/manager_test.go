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
		*NewSource(doc, "source1.php", SFILEREG),
		*NewSource(doc, "source2.php", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	m := Container.GetManager()
	results, err := m.buildStatus()

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
		*NewSource(doc, "source1.php", SFILEREG),
		*NewSource(doc, "source2.php", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	content := []byte("whatever")
	err = ioutil.WriteFile(util.GetAbsPath("source1.php"), content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	m := Container.GetManager()
	results, err := m.buildStatus()

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
		*NewSource(doc, "source1.php", SFILEREG),
		*NewSource(doc, "source2.php", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	err = os.Remove(util.GetAbsPath("source1.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	m := Container.GetManager()
	results, err := m.buildStatus()

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

func TestBuildStatusWithRegexpDescribingAFolder(t *testing.T) {
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
		*NewSource(doc, "test1", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	m := Container.GetManager()
	results, err := m.buildStatus()

	assert.NoError(t, err, "Must returns no errors")

	sort.Strings((*results)[0].Status[INONE])

	assert.Len(t, *results, 1, "One result must be returned")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[0].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "test1", (*results)[0].Source.Identifier, "Wrong source identifier returned")
	assert.Len(t, (*results)[0].Status[INONE], 3, "Three file items must be found")
	assert.Equal(t, []string{"test1/source1.php", "test1/source2.php", "test1/source3.php"}, (*results)[0].Status[INONE], "Three file items must be found")
}

func TestBuildStatusWithRegexpDescribingAFolderAndAddedFile(t *testing.T) {
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
		*NewSource(doc, "test1", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	createSourceFile([]byte("test"), "test1/source4.php")

	m := Container.GetManager()
	results, err := m.buildStatus()

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

func TestBuildStatusWithRegexpDescribingAFolderAndFolderDeleted(t *testing.T) {
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
		*NewSource(doc, "test1", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	err = os.RemoveAll(util.GetAbsPath("test1"))

	if err != nil {
		logrus.Fatal(err)
	}

	m := Container.GetManager()
	results, err := m.buildStatus()

	assert.NoError(t, err, "Must returns no errors")

	sort.Strings((*results)[0].Status[IDELETED])

	assert.Len(t, *results, 1, "One result must be returned")

	assert.Equal(t, "doc_file_to_track.txt", (*results)[0].Doc.Identifier, "Wrong doc identifier returned")
	assert.Equal(t, "test1", (*results)[0].Source.Identifier, "Wrong source identifier returned")

	assert.Len(t, (*results)[0].Status[IDELETED], 1, "One file item must be found")
	assert.Equal(t, []string{"test1/source1.php"}, (*results)[0].Status[IDELETED], "One file item must be found")
}

func retrieveItems(identifiers []string) map[string]*[]Item {
	items := map[string]*[]Item{}

	for _, identifier := range identifiers {
		var id string
		err := Container.GetDatabase().QueryRow("select id from sources where identifier = ?", identifier).Scan(&id)

		if err != nil {
			logrus.Warn(err)
		}

		source := Source{ID: id}
		items[identifier], err = Container.GetItemRepository().ListFromSource(&source)

		if err != nil {
			logrus.Warn(err)
		}
	}

	return items
}

func TestUpdateFingerprints(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "source1.php", SFILEREG),
		*NewSource(doc, "source2.php", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	before := retrieveItems([]string{"source1.php", "source2.php"})

	err = ioutil.WriteFile(util.GetAbsPath("source1.php"), []byte("<?php echo 'Hello world !';"), 0644)

	if err != nil {
		logrus.Fatal(err)
	}

	err = Container.GetManager().UpdateFingerprints()

	assert.NoError(t, err, "Must produces no errors")

	after := retrieveItems([]string{"source1.php", "source2.php"})

	assert.True(t, (*after["source1.php"])[0].UpdatedAt.After((*before["source1.php"])[0].UpdatedAt), "Must changes updated date")

	assert.True(t, (*after["source2.php"])[0].UpdatedAt.After((*before["source1.php"])[0].UpdatedAt), "Must changes updated date")
}

func TestDeleteItems(t *testing.T) {
	createMocks()
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "source1.php", SFILEREG),
		*NewSource(doc, "source2.php", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	before := retrieveItems([]string{"source1.php", "source2.php"})

	err = os.Remove(util.GetAbsPath("source1.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	err = Container.GetManager().UpdateFingerprints()

	assert.NoError(t, err, "Must produces no errors")

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
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "source1.php", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	doc = NewDoc("doc_file_to_track_2.txt", DFILE)

	sources = []Source{
		*NewSource(doc, "source2.php", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	before := retrieveItems([]string{"source1.php", "source2.php"})

	err = os.Remove(util.GetAbsPath("source1.php"))

	if err != nil {
		logrus.Fatal(err)
	}

	err = Container.GetManager().UpdateFingerprints()

	assert.NoError(t, err, "Must produces no errors")

	sourceRows, err := Container.GetDatabase().Query("select s.id from sources s where identifier = ?", "source1.php")

	if err != nil {
		logrus.Fatal(err)
	}

	docRows, err := Container.GetDatabase().Query("select d.id from docs d where identifier = ?", "doc_file_to_track.txt")

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
	err := Initialize()

	if err != nil {
		logrus.Fatal(err)
	}

	createSubTestDirectory("test1")
	createSourceFilesInPath("test1")

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "test1", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

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

	err = Container.GetManager().UpdateFingerprints()

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

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "test1", SFILEREG),
		*NewSource(doc, "test1/source5.php", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

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

	doc = NewDoc("doc_file_to_track_2.txt", DFILE)

	sources = []Source{
		*NewSource(doc, "test1/source6.php", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	createSourceFile([]byte("test"), "test1/source11.php")

	itemStatus, changesOccured, err := Container.GetManager().GetItemStatus()

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

	doc := NewDoc("doc_file_to_track.txt", DFILE)

	sources := []Source{
		*NewSource(doc, "source1.php", SFILEREG),
	}

	err = Container.GetConfigRepository().CreateFromDocAndSources(doc, &sources)

	if err != nil {
		logrus.Fatal(err)
	}

	itemStatus, changesOccured, err := Container.GetManager().GetItemStatus()

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
