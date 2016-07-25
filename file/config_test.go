package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIdentifierWithAFile(t *testing.T) {
	createMocks()

	value, category := ParseIdentifier("source1.php")
	assert.Equal(t, "source1.php", value, "Must return source file")
	assert.EqualValues(t, SFILE, category, "Must return file category")
}

func TestParseIdentifierWithAFolder(t *testing.T) {
	createMocks()
	createSubTestDirectory("test1")

	value, category := ParseIdentifier("test1")
	assert.Equal(t, "test1", value, "Must return source folder")
	assert.EqualValues(t, SFOLDER, category, "Must return folder category")
}

func TestParseIdentifierWithAnError(t *testing.T) {
	createMocks()

	value, category := ParseIdentifier("whatever")
	assert.Equal(t, "whatever", value, "Must return identifier")
	assert.EqualValues(t, SERROR, category, "Must return an error")
}
