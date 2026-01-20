package converter

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

type headingResult struct {
	lines map[int]headingInfo
	err   error
}

type listResult struct {
	lines map[int]listItemInfo
	err   error
}

func Convert(markdown []byte) (string, error) {
	doc := goldmark.New().Parser().Parse(text.NewReader(markdown))

	headingChan := make(chan headingResult)
	listChan := make(chan listResult)

	go func() {
		lines, err := extractHeadings(doc, markdown)
		headingChan <- headingResult{lines: lines, err: err}
	}()
	go func() {
		lines, err := extractListItems(doc, markdown)
		listChan <- listResult{lines: lines, err: err}
	}()

	headingRes := <-headingChan
	listRes := <-listChan

	if headingRes.err != nil {
		return "", headingRes.err
	}
	if listRes.err != nil {
		return "", listRes.err
	}

	return buildOutput(markdown, headingRes.lines, listRes.lines), nil
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
