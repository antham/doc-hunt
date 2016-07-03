package file

import (
	"fmt"
	"time"
)

// Update all entries in items database
func Update() error {
	err := updateItems()

	if err != nil {
		return err
	}

	err = updateItemsFingeprint()

	if err != nil {
		return err
	}

	return nil
}

// deleteItems remove items from their identifiers
func deleteItems(identifiers *[]string) error {
	var identifierQuery string

	for i, identifier := range *identifiers {
		identifierQuery += `"` + identifier + `"`

		if len(*identifiers)-1 != i {
			identifierQuery += ","
		}
	}

	_, err := db.Exec(fmt.Sprintf("delete from items where identifier in (%s)", identifierQuery))

	if err != nil {
		return err
	}

	_, err = db.Exec("delete from sources where id not in (select source_id from items);")

	if err != nil {
		return err
	}

	_, err = db.Exec("delete from docs where id not in (select doc_id from sources);")

	if err != nil {
		return err
	}

	return nil
}

// updateItems add missing occurence in database or removes those which disappeared
func updateItems() error {
	deleted := map[string]bool{}
	added := []Item{}

	for _, result := range *BuildStatus() {
		for _, filename := range result.Status[IDELETED] {
			deleted[filename] = true
		}

		if _, ok := result.Status[IADDED]; ok == true {
			items := result.Status[IADDED]

			for _, item := range *NewItems(&items, &result.Source) {
				added = append(added, item)
			}
		}
	}

	extractDeletedFiles := func(filenames *map[string]bool) *[]string {
		results := make([]string, len(*filenames))

		for filename := range *filenames {
			results = append(results, filename)
		}

		return &results
	}

	InsertItems(&added)
	err := deleteItems(extractDeletedFiles(&deleted))

	return err
}

// updateItemsFingeprint update file check sum if source file content changed
func updateItemsFingeprint() error {
	rows, err := db.Query("select distinct(identifier) from items")

	if err != nil {
		return err
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
			return err
		}

		_, err = db.Exec("update items set fingerprint = ?, updated_at = ? where identifier = ?", fingerprint, time.Now(), identifier)

		if err != nil {
			return err
		}
	}

	return nil
}
