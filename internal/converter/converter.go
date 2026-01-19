package converter

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func Convert(markdown []byte) (string, error) {
	doc := goldmark.New().Parser().Parse(text.NewReader(markdown))

	headingChan := make(chan map[int]headingInfo)
	listChan := make(chan map[int]listItemInfo)

	go func() { headingChan <- extractHeadings(doc, markdown) }()
	go func() { listChan <- extractListItems(doc, markdown) }()

	headingLines := <-headingChan
	listLines := <-listChan

	return buildOutput(markdown, headingLines, listLines), nil
}

func buildOutput(
	markdown []byte,
	headingLines map[int]headingInfo,
	listLines map[int]listItemInfo,
) string {
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
