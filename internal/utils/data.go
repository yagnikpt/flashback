package utils

import "strings"

func UniqueStrings(input []string) []string {
	seen := make(map[string]bool)
	var out []string

	for _, v := range input {
		lowerCase := strings.ToLower(v)
		if !seen[lowerCase] {
			seen[lowerCase] = true
			out = append(out, v)
		}
	}
	return out
}
