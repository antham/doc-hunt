package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateFingerprint(t *testing.T) {
	createMocks()

	fingerprint, err := calculateFingerprint("/tmp/doc-hunt/doc_file_to_track.txt")

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "8e2586d4d6c168565389214a17426a60f4bce67c", fingerprint, "Must calculate file fingerprint")
}

func TestCalculateFingerprintOfUnexistingFile(t *testing.T) {
	_, err := calculateFingerprint("/tmp/doc-hunt/whatever")

	assert.Error(t, err, "Must return an error")
}

func TestHasChanged(t *testing.T) {
	assert.True(t, hasChanged("8e2586d4d6c168565389214a17426a60f4bce67c", "444486d4d6c168565389214a17426a60f4bce67c"), "Must return true if file have a different checksum")
}

func TestHasChangedWithSameChecksum(t *testing.T) {
	assert.False(t, hasChanged("8e2586d4d6c168565389214a17426a60f4bce67c", "8e2586d4d6c168565389214a17426a60f4bce67c"), "Must return false if file have the same checksum")
}
