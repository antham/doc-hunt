package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAbsPath(t *testing.T) {
	AppPath = "/tmp"

	assert.Equal(t, "/tmp/file", GetAbsPath("file"), "Must return an absolute path")
}

func TestTrimAbsBasePath(t *testing.T) {
	AppPath = "/tmp"

	assert.Equal(t, "test", TrimAbsBasePath("/tmp/test"), "Must return relative path")
}

func TestGetFolderPathAddTrailingSeparator(t *testing.T) {
	assert.Equal(t, "test/", GetFolderPath("test"), "Must add a trailing seperator")
}

func TestGetFolderPathAddTrailingSeparatorWithMultipleTrailingSeparator(t *testing.T) {
	assert.Equal(t, "test/", GetFolderPath("test////////////"), "Must add a trailing seperator")
}
