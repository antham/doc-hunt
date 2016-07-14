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
func BuildStatus() (*[]Result, error) {
	results := []Result{}

	configs, err := ListConfig()

	if err != nil {
		return nil, err
	}

	for _, config := range *configs {
		for _, source := range config.Sources {
			result := Result{
				Status: map[ItemStatus][]string{},
			}

			result.Doc = config.Doc
			result.Source = source

			switch source.Category {
			case SFILE:
				status, err := retrieveFileItem(&source)

				if err != nil {
					return nil, err
				}

				result.Status[status] = append(result.Status[status], source.Identifier)
			case SFOLDER:
				folderItems, err := retrieveFolderItems(&source)

				if err != nil {
					return nil, err
				}

				for identifier, status := range *folderItems {
					result.Status[status] = append(result.Status[status], identifier)
				}
			}

			results = append(results, result)
		}
	}

	return &results, nil
}

func retrieveFolderItems(source *Source) (*map[string]ItemStatus, error) {
	items := map[string]ItemStatus{}
	dbItems, err := getItems(source)

	if err != nil {
		return nil, err
	}

	for _, item := range *dbItems {
		items[item.Identifier] = getFileStatus(item.Identifier, item.Fingerprint)
	}

	files, err := util.ExtractFolderFiles(source.Identifier)

	switch err.(type) {
	case nil:
	case *os.PathError:
	default:
		return &items, err
	}

	for _, file := range *files {
		if _, ok := items[file]; !ok {
			items[file] = IADDED
		}
	}

	return &items, nil
}

func retrieveFileItem(source *Source) (ItemStatus, error) {
	items, err := getItems(source)

	if err != nil {
		return INONE, err
	}

	return getFileStatus((*items)[0].Identifier, (*items)[0].Fingerprint), nil
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
