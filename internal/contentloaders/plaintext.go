package contentloaders

import (
	"errors"
	"os"
)

// txt, md and other text files
func GetTextFileContent(path string, response chan<- string, errChan chan<- error) {
	if path == "" {
		errChan <- errors.New("path cannot be empty")
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		errChan <- errors.New("file does not exist")
		return
	}

	content, err := os.ReadFile(path)
	if err != nil {
		errChan <- err
		return
	}
	response <- string(content)
	// chunks, err := utils.SplitText(string(content))
	// if err != nil {
	// 	errChan <- err
	// 	return
	// }

	// for _, chunk := range chunks {
	// 	response <- chunk
	// }
}
