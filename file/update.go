package file

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

// UpdateSourcesFingeprint update file check sum if source file content changed
func UpdateSourcesFingeprint() {
	rows, err := db.Query("select distinct(path) from sources")

	if err != nil {
		logrus.Fatal(err)
	}

	paths := []string{}

	for rows.Next() {
		var path string

		rows.Scan(&path)

		paths = append(paths, path)
	}

	for _, path := range paths {
		fingerprint, err := calculateFingerprint(path)

		if err != nil {
			logrus.Fatal(err)
		}

		_, err = db.Exec("update sources set fingerprint = ?, updated_at = ? where path = ?", fingerprint, time.Now(), path)

		if err != nil {
			logrus.Fatal(err)
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
		logrus.Fatal(err)
	}

	_, err = db.Exec("delete from docs where id not in (select doc_id from sources);")

	if err != nil {
		logrus.Fatal(err)
	}
}

// UpdateFilenameSources update filename sources when locations changed
func UpdateFilenameSources(filenames *map[string]string) {
	for origFilename, updatedFilename := range *filenames {
		_, err := db.Exec("update sources set path = ? where path = ?", updatedFilename, origFilename)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}
