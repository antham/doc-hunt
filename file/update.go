package file

import (
	"fmt"
	"time"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

// UpdateItemsFingeprint update file check sum if source file content changed
func UpdateItemsFingeprint() {
	rows, err := db.Query("select distinct(identifier) from items")

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}

	identifiers := []string{}

	for rows.Next() {
		var identifier string

		rows.Scan(&identifier)

		identifiers = append(identifiers, identifier)
	}

	for _, identifier := range identifiers {
		fingerprint, err := calculateFingerprint(identifier)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}

		_, err = db.Exec("update items set fingerprint = ?, updated_at = ? where identifier = ?", fingerprint, time.Now(), identifier)

		if err != nil {
			ui.Error(err)

			util.ErrorExit()
		}
	}
}

// DeleteItems remove items from their identifiers
func DeleteItems(identifiers *[]string) {
	var identifierQuery string

	for i, identifier := range *identifiers {
		identifierQuery += `"` + identifier + `"`

		if len(*identifiers)-1 != i {
			identifierQuery += ","
		}
	}

	_, err := db.Exec(fmt.Sprintf("delete from items where identifier in (%s)", identifierQuery))

	if err != nil {
		ui.Error(err)

		util.ErrorExit()
	}

	_, err = db.Exec("delete from sources where id not in (select source_id from items);")

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
