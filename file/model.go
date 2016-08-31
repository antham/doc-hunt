package file

import (
	"time"

	"github.com/satori/go.uuid"
)

// SourceCategory represents a source category
type SourceCategory int

// SourceCategory categories
const (
	SERROR = iota
	SFILEREG
)

// Source represents a source that we want to follow changes
type Source struct {
	ID         string
	Category   SourceCategory
	Identifier string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DocID      string
}

// NewSource create new source
func NewSource(doc *Doc, identifier string, category SourceCategory) *Source {
	source := Source{
		uuid.NewV4().String(),
		category,
		identifier,
		time.Now(),
		time.Now(),
		doc.ID,
	}

	return &source
}

// Item represents an actual tracked source
type Item struct {
	ID          string
	Identifier  string
	Fingerprint string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SourceID    string
}

// NewItems instanciate new items from an identifier list and a source
func NewItems(identifiers *[]string, source *Source) (*[]Item, error) {
	items := []Item{}

	for _, identifier := range *identifiers {
		fingerprint, err := calculateFingerprint(identifier)

		if err != nil {
			return nil, err
		}

		items = append(items, Item{
			uuid.NewV4().String(),
			identifier,
			fingerprint,
			time.Now(),
			time.Now(),
			source.ID,
		})
	}

	return &items, nil
}

// DocCategory represents a doc category
type DocCategory int

// DocCategory categories
const (
	DERROR = iota
	DFILE
	DURL
	DFOLDER
)

// Doc represents a document which as relationship with one or several source files
type Doc struct {
	ID         string
	Category   DocCategory
	Identifier string
	CreatedAt  time.Time
}

// NewDoc create a new doc file
func NewDoc(identifier string, category DocCategory) *Doc {
	return &Doc{
		uuid.NewV4().String(),
		category,
		identifier,
		time.Now(),
	}
}

// Setting represents a setting stored in database
type Setting struct {
	Name  string
	Value string
}

// NewSetting creates a new setting
func NewSetting(name string, value string) *Setting {
	return &Setting{
		name,
		value,
	}
}

// Config represents a config line
type Config struct {
	Doc     Doc
	Sources []Source
}

// Result represents what we get after comparison between database and actual files
type Result struct {
	Doc    Doc
	Source Source
	Status map[ItemStatus][]string
}
