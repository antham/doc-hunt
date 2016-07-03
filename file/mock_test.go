package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"

	"github.com/antham/doc-hunt/util"
)

func init() {
	util.AppPath = filepath.Clean(util.AppPath + "/../" + "test")
}

func removeTestDirectory() {
	err := os.RemoveAll(util.AppPath)

	if err != nil {
		logrus.Fatal(err)
	}
}

func createTestDirectory() {
	err := os.Mkdir(util.AppPath, 0777)

	if err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func createSubTestDirectory(path string) {
	err := os.Mkdir(fmt.Sprintf("%s%c%s", util.AppPath, filepath.Separator, path), 0777)

	if err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func createADocFile() {
	content := []byte("This is a doc file")
	err := ioutil.WriteFile(util.GetAbsPath("/doc_file_to_track.txt"), content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}
}

func createDocFiles() {
	for i := 1; i <= 10; i++ {
		content := []byte("This is a doc file")
		err := ioutil.WriteFile(util.GetAbsPath(fmt.Sprintf("doc_file_to_track_%d.txt", i)), content, 0644)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func createSourceFiles() {
	for i := 1; i <= 10; i++ {
		createSourceFile([]byte("<?php echo 'A source file';"), fmt.Sprintf("source%d.php", i))
	}
}

func createSourceFilesInPath(path string) {
	for i := 1; i <= 10; i++ {
		createSourceFile([]byte("<?php echo 'A source file';"), fmt.Sprintf("%s%csource%d.php", path, filepath.Separator, i))
	}
}

func createSourceFile(content []byte, filename string) {
	err := ioutil.WriteFile(util.GetAbsPath(filename), content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}
}

func createMocks() {
	removeTestDirectory()
	createTestDirectory()
	createADocFile()
	createSourceFiles()
}

func deleteDatabase() {
	os.Remove(util.GetAbsPath(dbName))
}
