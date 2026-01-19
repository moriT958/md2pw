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

type listItemInfo struct {
	level     int    // ネストレベル (1-3)
	isOrdered bool   // true: ordered (+), false: unordered (-)
	text      string // リストアイテムのテキスト
}

func Convert(markdown []byte) (string, error) {
	doc := goldmark.New().Parser().Parse(text.NewReader(markdown))
	headingLines := extractHeadings(doc, markdown)
	listLines := extractListItems(doc, markdown)
	return buildOutput(markdown, headingLines, listLines), nil
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

func extractListItems(doc ast.Node, markdown []byte) map[int]listItemInfo {
	listLines := make(map[int]listItemInfo)

	lineNumber := func(offset int) int {
		return bytes.Count(markdown[:offset], []byte("\n"))
	}

	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		li, ok := node.(*ast.ListItem)
		if !ok {
			return ast.WalkContinue, nil
		}

		// ネストレベルを計算（親の List ノードを数える）
		level := 0
		isOrdered := false
		for p := node.Parent(); p != nil; p = p.Parent() {
			if list, ok := p.(*ast.List); ok {
				level++
				if level == 1 {
					isOrdered = list.IsOrdered()
				}
			}
		}

		// 3レベルまで対応
		if level > 3 {
			level = 3
		}

		itemText := extractListItemText(li, markdown)
		line := lineNumber(li.FirstChild().Lines().At(0).Start)
		listLines[line] = listItemInfo{
			level:     level,
			isOrdered: isOrdered,
			text:      itemText,
		}
		return ast.WalkContinue, nil
	})

	return listLines
}

func extractListItemText(li *ast.ListItem, markdown []byte) string {
	var buf bytes.Buffer
	// ListItem の最初の子は通常 Paragraph
	if p := li.FirstChild(); p != nil {
		for child := p.FirstChild(); child != nil; child = child.NextSibling() {
			if t, ok := child.(*ast.Text); ok {
				buf.Write(t.Segment.Value(markdown))
			}
		}
	}
	return buf.String()
}

func buildOutput(markdown []byte, headingLines map[int]headingInfo, listLines map[int]listItemInfo) string {
	lines := strings.Split(string(markdown), "\n")
	var result bytes.Buffer

	for i, line := range lines {
		if h, ok := headingLines[i]; ok && h.level <= maxHeadingLevel {
			stars := strings.Repeat("*", h.level)
			result.WriteString(stars + " " + h.text)
		} else if li, ok := listLines[i]; ok {
			marker := strings.Repeat("-", li.level)
			if li.isOrdered {
				marker = strings.Repeat("+", li.level)
			}
			result.WriteString(marker + li.text)
		} else {
			result.WriteString(line)
		}
		if i < len(lines)-1 {
			result.WriteByte('\n')
		}
	}

	return result.String()
}
