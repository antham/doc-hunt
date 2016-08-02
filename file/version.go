package file

import (
	"fmt"
	"regexp"
	"strconv"
)

var appVersion = "1.0.0"

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

	res.Next()

	err = res.Scan(&version)

	return version, err
}

// HasMajorVersionEqualFrom check if major version in given version is equal to app version
func HasMajorVersionEqualFrom(ver string) (bool, error) {
	re := regexp.MustCompile(`^(\d)\.(\d)\.(\d)(?:(?:\-|\+).*)?`)

	verComp := re.FindStringSubmatch(ver)
	appVerComp := re.FindStringSubmatch(appVersion)

	if len(verComp) != 4 {
		return false, fmt.Errorf("Wrong version format : %s, must follows semver", ver)
	}

	appVerMajor, err := strconv.Atoi(appVerComp[1])

	if err != nil {
		return false, err
	}

	verMajor, err := strconv.Atoi(verComp[1])

	if err != nil {
		return false, err
	}

	return verMajor == appVerMajor, nil
}
