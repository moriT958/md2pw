package main

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Convert converts Markdown to PukiWiki format.
func Convert(markdown []byte) (string, error) {
	md := goldmark.New()
	reader := text.NewReader(markdown)
	doc := md.Parser().Parse(reader)

	// Build a map of heading line numbers to their info
	type headingInfo struct {
		level int
		text  string
	}
	headingLines := make(map[int]headingInfo)

	// Calculate line number from byte offset
	lineNumber := func(offset int) int {
		return bytes.Count(markdown[:offset], []byte("\n"))
	}

	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if h, ok := node.(*ast.Heading); ok {
			// Get heading text from children
			var textBuf bytes.Buffer
			for child := h.FirstChild(); child != nil; child = child.NextSibling() {
				if t, ok := child.(*ast.Text); ok {
					textBuf.Write(t.Segment.Value(markdown))
				}
			}
			// Find the line where this heading starts
			// Use the first child's segment to find the heading line
			if h.FirstChild() != nil {
				if t, ok := h.FirstChild().(*ast.Text); ok {
					line := lineNumber(t.Segment.Start)
					headingLines[line] = headingInfo{
						level: h.Level,
						text:  textBuf.String(),
					}
				}
			}
		}
		return ast.WalkContinue, nil
	})

	// Process line by line
	lines := strings.Split(string(markdown), "\n")
	var result bytes.Buffer
	headingPattern := regexp.MustCompile(`^#{1,6}\s+`)

	for i, line := range lines {
		if h, ok := headingLines[i]; ok && h.level <= 3 {
			stars := strings.Repeat("*", h.level)
			result.WriteString(stars + " " + h.text)
		} else if headingPattern.MatchString(line) {
			// This handles edge cases where AST didn't capture it
			result.WriteString(line)
		} else {
			result.WriteString(line)
		}
		if i < len(lines)-1 {
			result.WriteByte('\n')
		}
	}

	return result.String(), nil
}
