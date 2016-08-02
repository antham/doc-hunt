package file

import (
	"fmt"
	"regexp"
	"strconv"
)

var appVersion = "1.0.0"

// GetVersion return app version
func GetVersion() string {
	return appVersion
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
