package file

import (
	"os"

	"github.com/antham/doc-hunt/util"
)

//go:generate stringer -type=SourceStatus

// SourceStatus represents a source file status
type SourceStatus int

// Source files status
const (
	Failed SourceStatus = iota
	Untouched
	Updated
	Deleted
)

// FetchStatus retrieves sources file status
func FetchStatus() *[]Result {
	status := map[string]SourceStatus{}
	results := []Result{}

	for _, config := range *ListConfig() {
		result := Result{
			Status: map[SourceStatus][]string{},
		}
		result.Doc = config.Doc

		for _, source := range config.Sources {
			buildSourceStatus(&source, &status)

			result.Sources = append(result.Sources, source)
			result.Status[status[source.Path]] = append(result.Status[status[source.Path]], source.Path)
		}

		results = append(results, result)
	}

	return &results
}

func buildSourceStatus(source *Source, status *map[string]SourceStatus) {
	if _, ok := (*status)[(*source).Path]; ok {
		return
	}

	var fingerprint string

	filename := util.GetAbsPath(source.Path)

	_, err := os.Stat(filename)

	if err != nil {
		(*status)[(*source).Path] = Deleted

		return
	}

	fingerprint, err = calculateFingerprint(filename)

	if err != nil {
		(*status)[(*source).Path] = Failed

		return
	}

	if hasChanged(fingerprint, (*source).Fingerprint) {
		(*status)[(*source).Path] = Updated

		return
	}

	(*status)[(*source).Path] = Untouched
}
