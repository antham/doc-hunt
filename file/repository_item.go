package file

import (
	"database/sql"
	"fmt"
	"time"
)

// ItemRepository handle all actions availables on items table
type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository return a new ItemRepository instance
func NewItemRepository(db *sql.DB) ItemRepository {
	return ItemRepository{db: db}
}

// createFromItemList insert new items from an item list
func (i ItemRepository) createFromItemList(items *[]Item) error {
	for _, item := range *items {
		_, err := i.db.Exec("insert into items values (?,?,?,?,?,?)", item.ID, item.Identifier, item.Fingerprint, item.CreatedAt, item.UpdatedAt, item.SourceID)

		if err != nil {
			return err
		}
	}

	return nil
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

	_, err := i.db.Exec(fmt.Sprintf("delete from items where identifier in (%s)", identifierQuery))

	if err != nil {
		return err
	}

	_, err = i.db.Exec("delete from sources where id not in (select source_id from items);")

	if err != nil {
		return err
	}

	_, err = i.db.Exec("delete from docs where id not in (select doc_id from sources);")

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

	err = i.createFromItemList(&added)

	if err != nil {
		return err
	}

	return i.deleteByIdentifiers(extractDeletedFiles(&deleted))
}

// UpdateAllFingerprints update all items fingerprints with new value
// if file content changed
func (i ItemRepository) UpdateAllFingerprints() error {
	rows, err := i.db.Query("select distinct(identifier) from items")

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

		_, err = i.db.Exec("update items set fingerprint = ?, updated_at = ? where identifier = ?", fingerprint, time.Now(), identifier)

		if err != nil {
			return err
		}
	}

	return nil
}

// ListFromSource retrieves all items from a source
func (i ItemRepository) ListFromSource(source *Source) (*[]Item, error) {
	items := []Item{}

	rows, err := i.db.Query("select id, identifier, fingerprint, created_at, updated_at, source_id from items where source_id = ? order by identifier", source.ID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := Item{}

		err := rows.Scan(&item.ID, &item.Identifier, &item.Fingerprint, &item.CreatedAt, &item.UpdatedAt, &item.SourceID)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return &items, nil
}

// CreateFromIdentifiersAndSource insert new items from an identifier list and a source
func (i ItemRepository) CreateFromIdentifiersAndSource(identifiers *[]string, source *Source) error {
	items, err := NewItems(identifiers, source)

	if err != nil {
		return err
	}

	return i.createFromItemList(items)
}
