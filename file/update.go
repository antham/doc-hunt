package file

import (
	"fmt"
	"time"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

// UpdateSourcesFingeprint update file check sum if source file content changed
func UpdateSourcesFingeprint() {
	rows, err := db.Query("select distinct(path) from sources")

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}

	paths := []string{}

	for rows.Next() {
		var path string

		rows.Scan(&path)

		paths = append(paths, path)
	}

	for _, path := range paths {
		filename := dirApp + "/" + path
		fingerprint, err := calculateFingerprint(filename)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		_, err = db.Exec("update sources set fingerprint = ?, updated_at = ? where path = ?", fingerprint, time.Now(), path)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}
	}
}

// DeleteSources remove sources from their filenames
func DeleteSources(filenames *[]string) {
	var filenameQuery string

	for i, filename := range *filenames {
		filenameQuery += `"` + filename + `"`

		if len(*filenames)-1 != i {
			filenameQuery += ","
		}
	}

	_, err := db.Exec(fmt.Sprintf("delete from sources where path in (%s)", filenameQuery))

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}

	_, err = db.Exec("delete from docs where id not in (select doc_id from sources);")

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}
}

// UpdateFilenameSources update filename sources when locations changed
func UpdateFilenameSources(filenames *map[string]string) {
	for origFilename, updatedFilename := range *filenames {
		_, err := db.Exec("update sources set path = ? where path = ?", updatedFilename, origFilename)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}
	}
}
