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
