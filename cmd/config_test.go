package cmd

import (
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/antham/doc-hunt/util"
)

func init() {
	util.AppPath = util.AppPath + "/../" + "test"
}

func createTestDirectory() {
	err := os.RemoveAll(util.AppPath)

	if err != nil {
		logrus.Fatal(err)
	}

	err = os.Mkdir(util.AppPath, 0777)

	if err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func createDocFile() {
	content := []byte("")
	err := ioutil.WriteFile(util.GetAbsPath("doc_test_1"), content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}
}

func createSourceFiles() {
	content := []byte("")

	for _, name := range []string{util.GetAbsPath("source_test_1"), util.GetAbsPath("/source_test_2")} {
		err := ioutil.WriteFile(name, content, 0644)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}
