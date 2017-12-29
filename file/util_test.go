package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateFingerprint(t *testing.T) {
	createMocks()

	fingerprint, err := calculateFingerprint("doc_file_to_track.txt")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "8e2586d4d6c168565389214a17426a60f4bce67c", fingerprint, "Must calculate file fingerprint")
}

func TestCalculateFingerprintOfUnexistingFile(t *testing.T) {
	_, err := calculateFingerprint("whatever")

	assert.Error(t, err, "Must return an error")
}

func TestHasChanged(t *testing.T) {
	assert.True(t, hasChanged("8e2586d4d6c168565389214a17426a60f4bce67c", "444486d4d6c168565389214a17426a60f4bce67c"), "Must return true if file have a different checksum")
}

func TestHasChangedWithSameChecksum(t *testing.T) {
	assert.False(t, hasChanged("8e2586d4d6c168565389214a17426a60f4bce67c", "8e2586d4d6c168565389214a17426a60f4bce67c"), "Must return false if file have the same checksum")
}

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
