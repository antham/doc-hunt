package file

import (
	"fmt"
	"time"
)

// ItemRepository handle all actions availables on items table
type ItemRepository struct {
}

// NewItemRepository return a new ItemRepository instance
func NewItemRepository() ItemRepository {
	return ItemRepository{}
}

// deleteByIdentifiers all items from their identifiers
func (i ItemRepository) deleteByIdentifiers(identifiers *[]string) error {
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

	return err
}

// AppendOrRemove create new items from new files, delete those which were suppressed
func (i ItemRepository) AppendOrRemove() error {
	var err error

	deleted := map[string]bool{}
	added := []Item{}
	var status *[]Result

	status, err = BuildStatus()

	if err != nil {
		return err
	}

	for _, result := range *status {
		for _, filename := range result.Status[IDELETED] {
			deleted[filename] = true
		}

		if _, ok := result.Status[IADDED]; ok {
			itemsAdded := result.Status[IADDED]
			var items *[]Item

			items, err = NewItems(&itemsAdded, &result.Source)

			if err != nil {
				return err
			}

			added = append(added, *items...)
		}
	}

	extractDeletedFiles := func(filenames *map[string]bool) *[]string {
		results := make([]string, len(*filenames))

		for filename := range *filenames {
			results = append(results, filename)
		}

		return &results
	}

	err = InsertItems(&added)

	if err != nil {
		return err
	}

	return i.deleteByIdentifiers(extractDeletedFiles(&deleted))
}

// UpdateAllFingerprints update all items fingerprints with new value
// if file content changed
func (i ItemRepository) UpdateAllFingerprints() error {
	rows, err := db.Query("select distinct(identifier) from items")

	if err != nil {
		return err
	}

	identifiers := []string{}

	for rows.Next() {
		var identifier string

		err = rows.Scan(&identifier)

		if err != nil {
			return err
		}

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
