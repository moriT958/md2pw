package converter

import (
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

type codeblockResult struct {
	lines map[int]codeblockLineInfo
	err   error
}

type boldResult struct {
	bolds []boldInfo
	err   error
}

func Convert(markdown []byte) (string, error) {
	doc := goldmark.New().Parser().Parse(text.NewReader(markdown))

	headingChan := make(chan headingResult)
	listChan := make(chan listResult)
	codeblockChan := make(chan codeblockResult)
	boldChan := make(chan boldResult)

	go func() {
		lines, err := extractHeadings(doc, markdown)
		headingChan <- headingResult{lines: lines, err: err}
	}()
	go func() {
		lines, err := extractListItems(doc, markdown)
		listChan <- listResult{lines: lines, err: err}
	}()
	go func() {
		lines, err := extractCodeblocks(doc, markdown)
		codeblockChan <- codeblockResult{lines: lines, err: err}
	}()
	go func() {
		bolds, err := extractBolds(doc, markdown)
		boldChan <- boldResult{bolds: bolds, err: err}
	}()

	headingRes := <-headingChan
	listRes := <-listChan
	codeblockRes := <-codeblockChan
	boldRes := <-boldChan

	if headingRes.err != nil {
		return "", headingRes.err
	}
	if listRes.err != nil {
		return "", listRes.err
	}
	if codeblockRes.err != nil {
		return "", codeblockRes.err
	}
	if boldRes.err != nil {
		return "", boldRes.err
	}

	return buildOutput(markdown, headingRes.lines, listRes.lines, codeblockRes.lines, boldRes.bolds), nil
}

func buildOutput(
	markdown []byte,
	headingLines map[int]headingInfo,
	listLines map[int]listItemInfo,
	codeblockLines map[int]codeblockLineInfo,
	bolds []boldInfo,
) string {
	lines := strings.Split(string(markdown), "\n")
	var outputLines []string

	for i, line := range lines {
		if cb, ok := codeblockLines[i]; ok {
			if cb.isFence {
				continue // fence行をスキップ
			}
			outputLines = append(outputLines, cb.content)
		} else if h, ok := headingLines[i]; ok && h.level <= maxHeadingLevel {
			stars := strings.Repeat("*", h.level)
			outputLines = append(outputLines, stars+" "+h.text)
		} else if li, ok := listLines[i]; ok {
			marker := strings.Repeat("-", li.level)
			if li.isOrdered {
				marker = strings.Repeat("+", li.level)
			}
			outputLines = append(outputLines, marker+li.text)
		} else {
			outputLines = append(outputLines, line)
		}
	}

	output := strings.Join(outputLines, "\n")

	// Apply bold replacements
	for _, b := range bolds {
		output = strings.Replace(output, b.originalText, b.convertedText, 1)
	}

	return output
}
