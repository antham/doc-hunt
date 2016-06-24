package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAbsPath(t *testing.T) {
	AppPath = "/tmp"

	assert.Equal(t, "/tmp/file", GetAbsPath("file"), "Must return an absolute path")
}
