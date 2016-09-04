package file

import (
	"os"
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

// Manager handle all logic needed by command module
type Manager struct {
	cgfRepo                 *ConfigRepository
	itemRepo                *ItemRepository
	getAbsPath              func(string) string
	extractFilesMatchingReg func(string) (*[]string, error)
	stat                    func(string) (os.FileInfo, error)
}

// NewManager instanciate a manager
func NewManager(
	cgfRepo *ConfigRepository,
	itemRepo *ItemRepository,
	getAbsPath func(string) string,
	extractFilesMatchingReg func(string) (*[]string, error),
	stat func(string) (os.FileInfo, error),
) Manager {
	return Manager{
		cgfRepo:                 cgfRepo,
		itemRepo:                itemRepo,
		getAbsPath:              getAbsPath,
		extractFilesMatchingReg: extractFilesMatchingReg,
		stat: stat,
	}
}

// UpdateFingerprints update all entries in items database
func (m Manager) UpdateFingerprints() error {
	status, err := m.buildStatus()

	if err != nil {
		return err
	}

	err = m.itemRepo.AppendOrRemove(status)

	if err != nil {
		return err
	}

	return m.itemRepo.UpdateAllFingerprints()
}

// GetItemStatus retrieve source file status
func (m Manager) GetItemStatus() (*map[Doc]map[ItemStatus]map[string]bool, bool, error) {
	filenamePerStatus := map[Doc]map[ItemStatus]map[string]bool{}
	changesOccured := false

	results, err := m.buildStatus()

	if err != nil {
		return nil, false, err
	}

	for _, result := range *results {
		if filenamePerStatus[result.Doc] == nil {
			filenamePerStatus[result.Doc] = map[ItemStatus]map[string]bool{
				IADDED:   map[string]bool{},
				IUPDATED: map[string]bool{},
				IDELETED: map[string]bool{},
				IFAILED:  map[string]bool{},
				INONE:    map[string]bool{},
			}
		}

		for status, filenames := range result.Status {
			for _, filename := range filenames {
				if status == IDELETED || status == IUPDATED || status == IADDED || status == IFAILED {
					changesOccured = true
				}

				filenamePerStatus[result.Doc][status][filename] = true
			}
		}
	}

	return &filenamePerStatus, changesOccured, nil
}

// BuildStatus retrieves sources file status
func (m Manager) buildStatus() (*[]Result, error) {
	results := []Result{}

	configs, err := m.cgfRepo.List()

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
				items, err := m.findItems(&source)

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

func (m Manager) findItems(source *Source) (*map[string]ItemStatus, error) {
	items := map[string]ItemStatus{}

	dbItems, err := m.itemRepo.ListFromSource(source)

	if err != nil {
		return nil, err
	}

	for _, item := range *dbItems {
		items[item.Identifier] = m.getFileStatus(item.Identifier, item.Fingerprint)
	}

	files, err := m.extractFilesMatchingReg(source.Identifier)

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

func (m Manager) getFileStatus(path string, origFingerprint string) ItemStatus {
	filename := m.getAbsPath(path)

	_, err := m.stat(filename)

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
