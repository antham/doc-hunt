package file

import (
	"os"

	"github.com/antham/doc-hunt/util"
)

//go:generate stringer -type=ItemStatus

// ItemStatus represents a source file status
type ItemStatus int

// Source files status
const (
	INONE ItemStatus = iota
	IFAILED
	IADDED
	IUPDATED
	IDELETED
)

var statusCache map[string]ItemStatus

// BuildStatus retrieves sources file status
func BuildStatus() *[]Result {
	results := []Result{}

	for _, config := range *ListConfig() {
		for _, source := range config.Sources {
			result := Result{
				Status: map[ItemStatus][]string{},
			}

			result.Doc = config.Doc
			result.Source = source

			switch source.Category {
			case SFILE:
				status := retrieveFileItem(&source)
				result.Status[status] = append(result.Status[status], source.Identifier)
			case SFOLDER:
				for identifier, status := range retrieveFolderItems(&source) {
					result.Status[status] = append(result.Status[status], identifier)
				}
			}

			results = append(results, result)
		}
	}

	return &results
}

func retrieveFolderItems(source *Source) map[string]ItemStatus {
	items := map[string]ItemStatus{}

	for _, item := range *getItems(source) {
		items[item.Identifier] = getFileStatus(item.Identifier, item.Fingerprint)
	}

	for _, file := range *util.ExtractFolderFiles(source.Identifier) {
		if _, ok := items[file]; ok != true {
			items[file] = IADDED
		}
	}

	return items
}

func retrieveFileItem(source *Source) ItemStatus {
	item := (*getItems(source))[0]

	return getFileStatus(item.Identifier, item.Fingerprint)
}

func getFileStatus(path string, origFingerprint string) ItemStatus {
	filename := util.GetAbsPath(path)

	_, err := os.Stat(filename)

	if err != nil {
		return IDELETED
	}

	fingerprint, err := calculateFingerprint(path)

	if err != nil {
		return IFAILED
	}

	if hasChanged(fingerprint, origFingerprint) {
		return IUPDATED
	}

	return INONE
}
