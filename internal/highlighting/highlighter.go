package highlighting

import (
	"nanodocker/dockerfile"
	"regexp"
	"strings"
)

const (
	ColorReset   = "\033[0m"
	ColorKeyword = "\033[1;32m"
	ColorComment = "\033[0;36m"
	ColorString  = "\033[0;34m"
	ColorCode    = "\033[0;33m"
)

var dockerfileKeywords = []string{
	"FROM", "RUN", "CMD", "LABEL", "MAINTAINER", "EXPOSE", "ENV",
	"ADD", "COPY", "ENTRYPOINT", "VOLUME", "USER", "WORKDIR",
	"ARG", "ONBUILD", "STOPSIGNAL", "HEALTHCHECK", "SHELL",
}

func StripAnsiCodes(text string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(text, "")
}
func HighlightDockerfile(content string) string {
	lines := strings.Split(content, "\n")
	var highlightedLines []string

	for _, line := range lines {
		highlightedLines = append(highlightedLines, HighlightLine(line))
	}
	return strings.Join(highlightedLines, "\n")
}
func HighlightLine(line string) string {
	dockerfile.InitializeMap()
	trimmedLine := strings.TrimLeft(line, " \t\n\r")

	if strings.HasPrefix(trimmedLine, "#") {
		return ColorComment + line + ColorReset
	}

	for _, keyword := range dockerfileKeywords {
		if strings.HasPrefix(strings.ToUpper(trimmedLine), keyword) {
			keyWordEnd := len(keyword)
			if len(trimmedLine) >= keyWordEnd {
				leadingSpaces := len(line) - len(trimmedLine)
				actualKeyword := line[leadingSpaces : leadingSpaces+keyWordEnd]
				rest := ""
				if len(line) > leadingSpaces+keyWordEnd {
					rest = line[leadingSpaces+keyWordEnd:]
				}
				var highlighted string
				if len(rest) > 0 {
					if rest[0] == ' ' {
						highlighted = line[:leadingSpaces] + ColorKeyword + actualKeyword + ColorReset
					} else {
						highlighted = line[:leadingSpaces] + actualKeyword
					}
				} else {
					highlighted = line[:leadingSpaces] + ColorKeyword + actualKeyword + ColorReset
				}

				rest = HighlightStrings(rest)

				return highlighted + rest
			}
		}
	}
	return HighlightStrings(line)
}
func HighlightStrings(text string) string {
	stringRegex := regexp.MustCompile(`("[^"]*"|'[^']*')`)

	return stringRegex.ReplaceAllStringFunc(text, func(match string) string {
		return ColorString + match + ColorReset
	})
}
