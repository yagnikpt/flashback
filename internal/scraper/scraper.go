package scraper

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	h2m "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
	"github.com/yagnik-patel-47/flashback/internal/utils"
)

func GetPageContent(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	body := doc.Find("body")
	scripts := body.Find("script")
	scripts.Remove()
	html, err := body.Html()
	if err != nil {
		return nil, err
	}
	markdown, err := h2m.ConvertString(html)
	if err != nil {
		return nil, err
	}

	chunks, err := utils.SplitText(markdown)
	if err != nil {
		return nil, err
	}

	return chunks, nil
}

func ExtractURLs(text string) []string {
	urlPattern := `(?i)\b((?:https?://|www\.|[a-z0-9.-]+\.[a-z]{2,4}/?)[^\s<>"]+|www\.[^\s<>"]+)`
	re := regexp.MustCompile(urlPattern)
	urls := re.FindAllString(text, -1)

	var cleanedURLs []string
	for _, url := range urls {
		url = strings.TrimRight(url, ".,;:!?)")
		if strings.HasPrefix(url, "www.") {
			url = "http://" + url
		}
		cleanedURLs = append(cleanedURLs, url)
	}

	return cleanedURLs
}
