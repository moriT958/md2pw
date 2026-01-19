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
			input:    []byte("# Title\n\nSome paragraph text.\n\n## Section\n\nPlain text."),
			expected: "* Title\n\nSome paragraph text.\n\n** Section\n\nPlain text.",
		},
		// Unordered List のテストケース
		{
			name:     "基本的な unordered list",
			input:    []byte("- item1\n- item2\n- item3"),
			expected: "-item1\n-item2\n-item3",
		},
		{
			name:     "ネストした unordered list",
			input:    []byte("- item1\n  - nested1\n  - nested2\n- item2"),
			expected: "-item1\n--nested1\n--nested2\n-item2",
		},
		{
			name:     "3レベルのネスト unordered list",
			input:    []byte("- level1\n  - level2\n    - level3"),
			expected: "-level1\n--level2\n---level3",
		},
		// Ordered List のテストケース
		{
			name:     "基本的な ordered list",
			input:    []byte("1. item1\n2. item2\n3. item3"),
			expected: "+item1\n+item2\n+item3",
		},
		{
			name:     "ネストした ordered list",
			input:    []byte("1. item1\n   1. nested1\n   2. nested2\n2. item2"),
			expected: "+item1\n++nested1\n++nested2\n+item2",
		},
		{
			name:     "3レベルのネスト ordered list",
			input:    []byte("1. level1\n   1. level2\n      1. level3"),
			expected: "+level1\n++level2\n+++level3",
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
