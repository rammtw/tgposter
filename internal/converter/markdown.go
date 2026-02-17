package converter

import (
	"fmt"
	"regexp"
	"strings"
)

var tgSpecialChars = `_[]()~` + "`" + `>#+-=|{}.!\`

func MarkdownToTelegram(md string) string {
	lines := strings.Split(md, "\n")
	var result []string

	inCodeBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			result = append(result, line)
			continue
		}

		if inCodeBlock {
			line = escapeCodeContent(line)
			result = append(result, line)
			continue
		}

		// Заголовки → bold
		var headerRe = regexp.MustCompile(`^#{1,6}\s+(.+)$`)
		if m := headerRe.FindStringSubmatch(line); m != nil {
			escaped := escapeText(m[1]) // m[1] — текст заголовка
			result = append(result, "*"+escaped+"*")
			continue
		}

		line = convertFormattedLine(line)
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func convertFormattedLine(line string) string {
	type segment struct {
		placeholder string
		replacement string
	}
	var segments []segment
	counter := 0

	genPH := func() string {
		counter++
		return fmt.Sprintf("\x00PH%d\x00", counter)
	}

	codeRe := regexp.MustCompile("`([^`]+)`")
	line = codeRe.ReplaceAllStringFunc(line, func(m string) string {
		inner := codeRe.FindStringSubmatch(m)[2]
		ph := genPH()
		segments = append(segments, segment{ph, "`" + escapeCodeContent(inner) + "`"})
		return ph
	})

	linkRe := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	line = linkRe.ReplaceAllStringFunc(line, func(m string) string {
		parts := linkRe.FindStringSubmatch(m)
		text := escapeText(parts[1])
		url := parts[2]
		ph := genPH()
		segments = append(segments, segment{ph, "[" + text + "](" + url + ")"})
		return ph
	})

	boldRe := regexp.MustCompile(`\*\*(.+?)\*\*`)
	line = boldRe.ReplaceAllStringFunc(line, func(m string) string {
		inner := boldRe.FindStringSubmatch(m)[1]
		ph := genPH()
		segments = append(segments, segment{ph, "*" + escapeText(inner) + "*"})
		return ph
	})

	italicRe := regexp.MustCompile(`\*(.+?)\*`)
	line = italicRe.ReplaceAllStringFunc(line, func(m string) string {
		inner := italicRe.FindStringSubmatch(m)[1]
		ph := genPH()
		segments = append(segments, segment{ph, "_" + escapeText(inner) + "_"})
		return ph
	})

	strikeRe := regexp.MustCompile(`~~(.+?)~~`)
	line = strikeRe.ReplaceAllStringFunc(line, func(m string) string {
		inner := strikeRe.FindStringSubmatch(m)[1]
		ph := genPH()
		segments = append(segments, segment{ph, "~" + escapeText(inner) + "~"})
		return ph
	})

	line = escapeText(line)

	for _, seg := range segments {
		line = strings.Replace(line, seg.placeholder, seg.replacement, 1)
	}

	return line
}

func escapeText(s string) string {
	var b strings.Builder
	for _, r := range s {
		if strings.ContainsRune(tgSpecialChars, r) {
			b.WriteRune('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

func escapeCodeContent(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "`", "\\`")
	return s
}
