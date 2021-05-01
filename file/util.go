package file

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"

	"github.com/dlclark/regexp2"

	"github.com/antham/doc-hunt/util"
)

func calculateFingerprint(filepath string) (string, error) {
	data, err := ioutil.ReadFile(util.GetAbsPath(filepath))

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1.Sum(data)), nil
}

func hasChanged(actualFileSum string, backupFileSum string) bool {
	return actualFileSum != backupFileSum[:]
}

// ParseIdentifier extract identifier and category from string
func ParseIdentifier(value string) (string, SourceCategory) {
	_, err := regexp2.Compile(value, regexp2.ExplicitCapture)

	if err == nil {
		return value, SFILEREG
	}

	return value, SERROR
}
