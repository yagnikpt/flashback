package contentloaders

import (
	"errors"
	"os"

	"github.com/yagnikpt/flashback/internal/utils"
)

// txt, md and other text files
func GetTextContent(path string) ([]string, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("file does not exist")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	chunks, err := utils.SplitText(string(content))
	if err != nil {
		return nil, err
	}

	return chunks, nil
}
