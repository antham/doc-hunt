package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIdentifierWithARegexp(t *testing.T) {
	value, category := ParseIdentifier("test.*")
	assert.Equal(t, "test.*", value, "Must return identifier")
	assert.EqualValues(t, SFILEREG, category, "Must return file regexp category")
}

func TestParseIdentifierWithAnError(t *testing.T) {
	createMocks()

	value, category := ParseIdentifier("test****")
	assert.Equal(t, "test****", value, "Must return identifier")
	assert.EqualValues(t, SERROR, category, "Must return an error")
}
