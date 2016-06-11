package file

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
)

func calculateFingerprint(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha1.Sum(data)), nil
}

func hasChanged(actualFileSum string, backupFileSum string) bool {
	return actualFileSum != backupFileSum[:]
}
