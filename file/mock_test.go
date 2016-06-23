package file

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
)

func init() {
	currentDir, dirErr := os.Getwd()

	if dirErr != nil {
		logrus.Fatal(dirErr)
	}

	dirApp = currentDir + "/../" + "test"
}

func removeTestDirectory() {
	err := os.RemoveAll(dirApp)

	if err != nil {
		logrus.Fatal(err)
	}
}

func createTestDirectory() {
	err := os.Mkdir(dirApp, 0777)

	if err != nil && !os.IsExist(err) {
		logrus.Fatal(err)
	}
}

func createADocFile() {
	content := []byte("This is a doc file")
	err := ioutil.WriteFile(dirApp+"/doc_file_to_track.txt", content, 0644)

	if err != nil {
		logrus.Fatal(err)
	}
}

func createDocFiles() {
	for i := 1; i <= 10; i++ {
		content := []byte("This is a doc file")
		err := ioutil.WriteFile(fmt.Sprintf(dirApp+"/doc_file_to_track_%d.txt", i), content, 0644)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func createSourceFiles() {
	for i := 1; i <= 10; i++ {
		content := []byte("<?php echo 'A source file';")
		err := ioutil.WriteFile(fmt.Sprintf(dirApp+"/source%d.php", i), content, 0644)

		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func createMocks() {
	removeTestDirectory()
	createTestDirectory()
	createADocFile()
	createSourceFiles()
}

func deleteDatabase() {
	err := os.Remove(dbName)

	if err != nil {
		logrus.Fatal(err)
	}
}
