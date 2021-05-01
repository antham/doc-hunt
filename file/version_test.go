package file

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetAppVersion(t *testing.T) {
	assert.Equal(t, appVersion, GetAppVersion(), "Must return version")
}

func TestHasMajorVersionEqualFromWithSameMajorVersion(t *testing.T) {
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Error(err)
	}

	appVersion = "1.0.0"

	_, err = db.Exec("update version set id = '1.6.9'")

	if err != nil {
		logrus.Error(err)
	}

	isEqual, err := HasMajorVersionEqualFrom()

	assert.NoError(t, err, "Must return no errors")
	assert.True(t, isEqual, "Must be true, major versions are equals")
}

func TestHasMajorVersionEqualFromWithDifferentMajorVersion(t *testing.T) {
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Error(err)
	}

	appVersion = "1.0.0"

	_, err = db.Exec("update version set id = '2.0.0'")

	if err != nil {
		logrus.Error(err)
	}

	isEqual, err := HasMajorVersionEqualFrom()

	assert.EqualError(t, err, "Database version : 2.0.0 and app version : 1.0.0 don't have same major version", "Must return an error if version format are different")
	assert.False(t, isEqual, "Must be false, major versions are differents")
}

func TestHasMajorVersionEqualFromWithWrongVersionFormat(t *testing.T) {
	deleteDatabase()
	err := Initialize()

	if err != nil {
		logrus.Error(err)
	}

	appVersion = "1.0.0"

	_, err = db.Exec("update version set id = '2.0'")

	if err != nil {
		logrus.Error(err)
	}

	_, err = HasMajorVersionEqualFrom()

	assert.EqualError(t, err, "Wrong version format : 2.0, must follows semver", "Must return an error if version format is invalid")
}
