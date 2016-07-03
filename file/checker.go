package file

// GetItemStatus retrieve source file status
func GetItemStatus() *map[Doc]map[ItemStatus]map[string]bool {
	filenamePerStatus := map[Doc]map[ItemStatus]map[string]bool{}

	for _, result := range *BuildStatus() {
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
				filenamePerStatus[result.Doc][status][filename] = true
			}
		}
	}

	return &filenamePerStatus
}
