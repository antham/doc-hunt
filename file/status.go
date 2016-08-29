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

// BuildStatus retrieves sources file status
func BuildStatus() (*[]Result, error) {
	results := []Result{}

	c := Container.GetConfigRepository()
	configs, err := c.List()

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
			case SFILEREG:
				items, err := findItems(&source)

				if err != nil {
					return nil, err
				}

				for identifier, status := range *items {
					result.Status[status] = append(result.Status[status], identifier)
				}
			}

			results = append(results, result)
		}
	}

	return &results, nil
}

func findItems(source *Source) (*map[string]ItemStatus, error) {
	items := map[string]ItemStatus{}

	r := Container.GetItemRepository()
	dbItems, err := r.ListFromSource(source)

	if err != nil {
		return nil, err
	}

	for _, item := range *dbItems {
		items[item.Identifier] = getFileStatus(item.Identifier, item.Fingerprint)
	}

	files, err := util.ExtractFilesMatchingReg(source.Identifier)

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
