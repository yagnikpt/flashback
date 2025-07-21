package contentloaders

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	h2m "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
	"github.com/yagnikpt/flashback/internal/utils"
	"google.golang.org/genai"
)

func GetWebpageContent(query string, url string, client *genai.Client, response chan<- string, errChan chan<- error) {
	url = strings.TrimRight(url, ".,;:!?)")
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	res, err := http.Get(url)
	if err != nil {
		errChan <- err
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		errChan <- fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		errChan <- err
		return
	}

	body := doc.Find("body")
	scripts := body.Find("script")
	scripts.Remove()
	html, err := body.Html()
	if err != nil {
		errChan <- err
		return
	}
	markdown, err := h2m.ConvertString(html)
	if err != nil {
		errChan <- err
		return
	}
	input := "text content: " + markdown + "\n\n" + "user query: " + query
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(string(utils.WebExtractionPrompt), genai.RoleUser),
	}
	result, err := client.Models.GenerateContent(
		context.Background(),
		"gemini-2.5-flash",
		genai.Text(input),
		config,
	)

	if err != nil {
		errChan <- err
		return
	}

	response <- result.Text()
}
