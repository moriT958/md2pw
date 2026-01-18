package converter

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

const maxHeadingLevel = 3

type headingInfo struct {
	level int
	text  string
}

func Convert(markdown []byte) (string, error) {
	doc := goldmark.New().Parser().Parse(text.NewReader(markdown))
	headingLines := extractHeadings(doc, markdown)
	return buildOutput(markdown, headingLines), nil
}

func extractHeadings(doc ast.Node, markdown []byte) map[int]headingInfo {
	headingLines := make(map[int]headingInfo)

	lineNumber := func(offset int) int {
		return bytes.Count(markdown[:offset], []byte("\n"))
	}

	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		h, ok := node.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		headingText := extractHeadingText(h, markdown)
		if t, ok := h.FirstChild().(*ast.Text); ok {
			line := lineNumber(t.Segment.Start)
			headingLines[line] = headingInfo{
				level: h.Level,
				text:  headingText,
			}
		}
		return ast.WalkContinue, nil
	})

	return headingLines
}

func extractHeadingText(h *ast.Heading, markdown []byte) string {
	var buf bytes.Buffer
	for child := h.FirstChild(); child != nil; child = child.NextSibling() {
		if t, ok := child.(*ast.Text); ok {
			buf.Write(t.Segment.Value(markdown))
		}
	}
	return buf.String()
}

func buildOutput(markdown []byte, headingLines map[int]headingInfo) string {
	lines := strings.Split(string(markdown), "\n")
	var result bytes.Buffer

	for i, line := range lines {
		if h, ok := headingLines[i]; ok && h.level <= maxHeadingLevel {
			stars := strings.Repeat("*", h.level)
			result.WriteString(stars + " " + h.text)
		} else {
			result.WriteString(line)
		}
		if i < len(lines)-1 {
			result.WriteByte('\n')
		}
	}

	return result.String()
}
