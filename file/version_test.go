package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	assert.Equal(t, appVersion, GetVersion(), "Must return version")
}

func TestHasMajorVersionEqualFromWithSameMajorVersion(t *testing.T) {
	appVersion = "1.0.0"

	isEqual, err := HasMajorVersionEqualFrom("1.6.9")

	assert.NoError(t, err, "Must return no errors")
	assert.True(t, isEqual, "Must be true, major versions are equals")
}

func TestHasMajorVersionEqualFromWithDifferentMajorVersion(t *testing.T) {
	appVersion = "1.0.0"

	isEqual, err := HasMajorVersionEqualFrom("2.0.0")

	assert.NoError(t, err, "Must return no errors")
	assert.False(t, isEqual, "Must be false, major versions are differents")
}

func TestHasMajorVersionEqualFromWithWrongVersionFormat(t *testing.T) {
	appVersion = "1.0.0"

	_, err := HasMajorVersionEqualFrom("2.0")

	assert.EqualError(t, err, "Wrong version format : 2.0, must follows semver", "Must return an error if version format is invalid")
}
