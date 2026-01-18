package converter

import (
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "H1",
			input:    []byte("# Heading 1"),
			expected: "* Heading 1",
		},
		{
			name:     "H2",
			input:    []byte("## Heading 2"),
			expected: "** Heading 2",
		},
		{
			name:     "H3",
			input:    []byte("### Heading 3"),
			expected: "*** Heading 3",
		},
		{
			name:     "複数種類の Heading がある場合も全て変換される",
			input:    []byte("# H1\n\n## H2\n\n### H3"),
			expected: "* H1\n\n** H2\n\n*** H3",
		},
		{
			name:     "H4 は変換されない",
			input:    []byte("# H1\n\n#### H4\n\n## H2"),
			expected: "* H1\n\n#### H4\n\n** H2",
		},
		{
			name:     "内容がそのまま保たれる",
			input:    []byte("# Title\n\nSome paragraph text.\n\n## Section\n\n- list item"),
			expected: "* Title\n\nSome paragraph text.\n\n** Section\n\n- list item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input)
			if err != nil {
				t.Fatalf("Convert returned error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
