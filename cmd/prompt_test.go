package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockTerminalReader struct {
	output string
	err    error
}

func (m *mockTerminalReader) Readline() (string, error) {
	if m.err != nil {
		return "", m.err
	}

	return m.output, nil
}

func (m *mockTerminalReader) Close() error {
	return nil
}

func TestBasePrompt(t *testing.T) {
	rl = &mockTerminalReader{"test", nil}

	callback := func(line string) error {
		return nil
	}

	assert.Equal(t, "test", basePrompt("A prompt", callback), "Must return data set on prompt")
}

func TestGlobalPromptDeleteOption(t *testing.T) {
	rl = &mockTerminalReader{"test", nil}

	filename := "test"
	deleted := map[string]bool{}
	moved := map[string]string{}

	err := globalPrompt(filename, &deleted, &moved)("d")

	assert.NoError(t, err, "Must throws no errors")
	assert.Equal(t, map[string]bool{"test": true}, deleted, "Must stores filename")
}

func TestGlobalPromptRenameOption(t *testing.T) {
	bckPrompt := basePrompt

	basePrompt = func(prompt string, callback checker) string {
		return ""
	}

	rl = &mockTerminalReader{"test", nil}

	filename := "test"
	deleted := map[string]bool{}
	moved := map[string]string{}

	err := globalPrompt(filename, &deleted, &moved)("r")

	assert.NoError(t, err, "Must throws no errors")

	basePrompt = bckPrompt
}

func TestGlobalPromptWrongOption(t *testing.T) {
	rl = &mockTerminalReader{"test", nil}

	filename := "test"
	deleted := map[string]bool{}
	moved := map[string]string{}

	err := globalPrompt(filename, &deleted, &moved)("w")

	assert.Error(t, err, "Must throws an error")
	assert.EqualError(t, err, `This action "w" : doesn't exist`, "Must throws an error when option is invalid")
}

func TestRenamePromptUnexistingFilename(t *testing.T) {
	rl = &mockTerminalReader{"test", nil}

	filename := "test"
	moved := map[string]string{}

	err := renamePrompt(filename, &moved)("/tmp/whatever")

	assert.Error(t, err, "Must throws an error")
	assert.EqualError(t, err, `File "/tmp/whatever" doesn't exist, please enter an existing filename`, "Must throws an error when filename doesn't exist")
}

func TestRenamePrompt(t *testing.T) {
	createTestDirectory()
	createSourceFiles()

	rl = &mockTerminalReader{"test", nil}

	filename := "test"
	moved := map[string]string{}

	err := renamePrompt(filename, &moved)("source_test_1")

	assert.NoError(t, err, "Must throws no error")
	assert.Equal(t, map[string]string{"test": "source_test_1"}, moved, "Must store original and renamed file")
}
