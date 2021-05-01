package file

import (
	"regexp"
)

var appVersion = "2.1.1"

// Version handle everything needed to compare versions,
// retrieve app version
type Version struct {
	settingRepo *SettingRepository
	mustCompile func(string) *regexp.Regexp
	atoi        func(string) (int, error)
	errorf      func(string, ...interface{}) error
}

// Get return app version
func (v Version) Get() string {
	return appVersion
}

// NewVersion instanciates a new version
func NewVersion(
	settingRepo *SettingRepository,
	mustCompile func(string) *regexp.Regexp,
	atoi func(string) (int, error),
	errorf func(string, ...interface{}) error,

) Version {
	return Version{
		settingRepo,
		mustCompile,
		atoi,
		errorf,
	}
}

// getDbVersion return app version recorded in database
func (v Version) getDbVersion() (string, error) {
	version, err := v.settingRepo.Get("version")

	if err != nil {
		return "", v.errorf("Can't retrieve database app version")
	}

	return version, nil
}

// HasMajorVersionEqual check if major version in given version is equal to app version
func (v Version) HasMajorVersionEqual() (bool, error) {
	re := v.mustCompile(`^(\d)\.(\d)\.(\d)(?:(?:\-|\+).*)?`)

	dbVer, err := v.getDbVersion()

	if err != nil {
		return false, err
	}

	dbVerComp := re.FindStringSubmatch(dbVer)
	appVerComp := re.FindStringSubmatch(appVersion)

	if len(dbVerComp) != 4 {
		return false, v.errorf("Wrong version format : %s, must follows semver", dbVer)
	}

	appVerMajor, err := v.atoi(appVerComp[1])

	if err != nil {
		return false, err
	}

	verMajor, err := v.atoi(dbVerComp[1])

	if err != nil {
		return false, err
	}

	if verMajor == appVerMajor {
		return true, nil
	}

	return false, v.errorf("Database version : %s and app version : %s don't have same major version", dbVer, appVersion)
}
