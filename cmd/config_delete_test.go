package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/file"
)

func TestParseConfigDelArgsWithArgumentNotANumber(t *testing.T) {
	configs := []file.Config{
		file.Config{},
		file.Config{},
		file.Config{},
		file.Config{},
	}

	_, err := parseConfigDelArgs(&configs, "1,2,3,a")

	assert.EqualError(t, err, "a is not a number", "Must return an error")
}

func TestParseConfigDelArgsWithArgumentNotInRange(t *testing.T) {
	configs := []file.Config{
		file.Config{},
		file.Config{},
		file.Config{},
	}

	_, err := parseConfigDelArgs(&configs, "3,4")

	assert.EqualError(t, err, "Value 3 is out of bounds", "Must return an error")
}

func TestParseConfigDelArgs(t *testing.T) {
	configs := []file.Config{
		file.Config{Doc: file.Doc{Identifier: "doc0.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc1.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc2.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc3.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc4.txt"}},
	}

	expected := &[]file.Config{
		file.Config{Doc: file.Doc{Identifier: "doc3.txt"}},
		file.Config{Doc: file.Doc{Identifier: "doc4.txt"}},
	}

	results, err := parseConfigDelArgs(&configs, "3,4")

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, expected, results, "Must return configs")
}
