package utils

import "github.com/tmc/langchaingo/textsplitter"

func SplitText(text string) ([]string, error) {
	split := textsplitter.NewRecursiveCharacter() // 512, 100
	split.Separators = []string{"\n", "\r\n", "\r"}

	chunks, err := split.SplitText(text)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}
