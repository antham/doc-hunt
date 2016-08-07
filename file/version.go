package file

import (
	"fmt"
	"regexp"
	"strconv"
)

var appVersion = "2.1.1"

// GetAppVersion return app version
func GetAppVersion() string {
	return appVersion
}

// getAppDbVersion return app version recorded in database
func getAppDbVersion() (string, error) {
	var version string

	res, err := db.Query("select id from version")

	if err != nil {
		return "", fmt.Errorf("Can't retrieve database app version")
	}

	defer func() {
		if e := res.Close(); e != nil {
			err = e
		}
	}()

	if err == nil {
		for res.Next() {
			err = res.Scan(&version)
		}
	}

	return version, err
}

// HasMajorVersionEqualFrom check if major version in given version is equal to app version
func HasMajorVersionEqualFrom() (bool, error) {
	re := regexp.MustCompile(`^(\d)\.(\d)\.(\d)(?:(?:\-|\+).*)?`)

	dbVer, err := getAppDbVersion()

	if err != nil {
		return false, err
	}

	dbVerComp := re.FindStringSubmatch(dbVer)
	appVerComp := re.FindStringSubmatch(appVersion)

	if len(dbVerComp) != 4 {
		return false, fmt.Errorf("Wrong version format : %s, must follows semver", dbVer)
	}

	appVerMajor, err := strconv.Atoi(appVerComp[1])

	if err != nil {
		return false, err
	}

	verMajor, err := strconv.Atoi(dbVerComp[1])

	if err != nil {
		return false, err
	}

	if verMajor == appVerMajor {
		return true, nil
	}

	return false, fmt.Errorf("Database version : %s and app version : %s don't have same major version", dbVer, appVersion)
}
