package cmd

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
