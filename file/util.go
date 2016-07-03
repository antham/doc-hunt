package file

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"

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
