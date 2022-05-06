package cmd

import (
	"bytes"
	"testing"

	// "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/doc-hunt/file"
)

func TestRenderDryRun(t *testing.T) {
	output := []byte{}
	out = bytes.NewBuffer(output)

	doc := file.Doc{}
	doc.Identifier = "doc-1.txt"

	sources := map[string]*[]string{
		"file/.*": {
			"file/source1.php",
			"file/source2.php",
		},
	}

	renderDryRun(&doc, &sources)

	assert.Regexp(t, `doc-1.txt`, out, "Must render document")
	assert.Regexp(t, `Files matching regexp "file/\.\*"`, out, "Must render original regexp")
	assert.Regexp(t, `=> file/source1.php`, out, "Must render source")
	assert.Regexp(t, `=> file/source2.php`, out, "Must render source")
}

func TestRenderDryRunWithEmptySources(t *testing.T) {
	output := []byte{}
	out = bytes.NewBuffer(output)

	doc := file.Doc{}
	doc.Identifier = "doc-1.txt"

	sources := map[string]*[]string{
		"file/.*": {},
	}

	renderDryRun(&doc, &sources)

	assert.Regexp(t, `doc-1.txt`, out, "Must render document")
	assert.Regexp(t, `Files matching regexp "file/\.\*"`, out, "Must render original request")
	assert.Regexp(t, `=> No files found`, out, "Must render an error message showing no sources are found")
}
