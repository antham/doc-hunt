package file

// GetItemStatus retrieve source file status
func GetItemStatus() (*map[Doc]map[ItemStatus]map[string]bool, bool, error) {
	filenamePerStatus := map[Doc]map[ItemStatus]map[string]bool{}
	changesOccured := false

	results, err := BuildStatus()

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
