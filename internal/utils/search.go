package utils

import (
	"regexp"
)

type SearchResult struct {
	Web   []string `json:"web"`
	Files []string `json:"files"`
}

func ExtractSearchTerms(input string) SearchResult {
	result := SearchResult{
		Web:   make([]string, 0),
		Files: make([]string, 0),
	}

	webRegex := regexp.MustCompile(`^web:(.+)$`)
	fileRegex := regexp.MustCompile(`^file:(.+)$`)

	webMatches := webRegex.FindAllStringSubmatch(input, -1)
	for _, match := range webMatches {
		if len(match) > 1 {
			result.Web = append(result.Web, match[1])
		}
	}

	fileMatches := fileRegex.FindAllStringSubmatch(input, -1)
	for _, match := range fileMatches {
		if len(match) > 1 {
			result.Files = append(result.Files, match[1])
		}
	}

	return result
}
