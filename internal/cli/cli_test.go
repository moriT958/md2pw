package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_StdinInput(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		args           []string
		expectedOutput string
		expectedCode   int
	}{
		{
			name:           "stdin with dash argument",
			input:          "# Heading",
			args:           []string{"md2pw", "-"},
			expectedOutput: "* Heading",
			expectedCode:   0,
		},
		{
			name:           "stdin with dash and output flag",
			input:          "- item1\n- item2",
			args:           []string{"md2pw", "-"},
			expectedOutput: "-item1\n-item2",
			expectedCode:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inStream := strings.NewReader(tt.input)
			outStream := &bytes.Buffer{}
			errStream := &bytes.Buffer{}

			c := New(inStream, outStream, errStream)
			code := c.Run(tt.args)

			if code != tt.expectedCode {
				t.Errorf("expected exit code %d, got %d. stderr: %s", tt.expectedCode, code, errStream.String())
			}
			if tt.expectedCode == 0 && outStream.String() != tt.expectedOutput {
				t.Errorf("expected output %q, got %q", tt.expectedOutput, outStream.String())
			}
		})
	}
}

func TestRun_FileInput(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(inputFile, []byte("# Test Heading"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedCode   int
	}{
		{
			name:           "file argument",
			args:           []string{"md2pw", inputFile},
			expectedOutput: "* Test Heading",
			expectedCode:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inStream := strings.NewReader("")
			outStream := &bytes.Buffer{}
			errStream := &bytes.Buffer{}

			c := New(inStream, outStream, errStream)
			code := c.Run(tt.args)

			if code != tt.expectedCode {
				t.Errorf("expected exit code %d, got %d. stderr: %s", tt.expectedCode, code, errStream.String())
			}
			if tt.expectedCode == 0 && outStream.String() != tt.expectedOutput {
				t.Errorf("expected output %q, got %q", tt.expectedOutput, outStream.String())
			}
		})
	}
}

func TestRun_FileArgumentPrioritizedOverStdin(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(inputFile, []byte("## File Content"), 0644); err != nil {
		t.Fatal(err)
	}

	// stdin has different content
	inStream := strings.NewReader("# Stdin Content")
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}

	c := New(inStream, outStream, errStream)
	code := c.Run([]string{"md2pw", inputFile})

	if code != 0 {
		t.Errorf("expected exit code 0, got %d. stderr: %s", code, errStream.String())
	}
	// File content should be used, not stdin
	expected := "** File Content"
	if outStream.String() != expected {
		t.Errorf("expected output %q, got %q", expected, outStream.String())
	}
}

func TestRun_OutputToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	inStream := strings.NewReader("# Test")
	outStream := &bytes.Buffer{}
	errStream := &bytes.Buffer{}

	c := New(inStream, outStream, errStream)
	code := c.Run([]string{"md2pw", "-o", outputFile, "-"})

	if code != 0 {
		t.Errorf("expected exit code 0, got %d. stderr: %s", code, errStream.String())
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	expected := "* Test"
	if string(content) != expected {
		t.Errorf("expected file content %q, got %q", expected, string(content))
	}
}

func TestRun_ErrorCases(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedCode int
	}{
		{
			name:         "nonexistent file",
			args:         []string{"md2pw", "/nonexistent/file.md"},
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inStream := strings.NewReader("")
			outStream := &bytes.Buffer{}
			errStream := &bytes.Buffer{}

			c := New(inStream, outStream, errStream)
			code := c.Run(tt.args)

			if code != tt.expectedCode {
				t.Errorf("expected exit code %d, got %d", tt.expectedCode, code)
			}
		})
	}
}
