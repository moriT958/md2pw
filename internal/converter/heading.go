package converter

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
)

const maxHeadingLevel = 3

type headingInfo struct {
	level int
	text  string
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
