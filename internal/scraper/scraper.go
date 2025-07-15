package scraper

import (
	"fmt"
	"net/http"
	"strings"

	h2m "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
	"github.com/yagnik-patel-47/flashback/internal/utils"
)

func GetPageContent(url string) ([]string, error) {
	url = strings.TrimRight(url, ".,;:!?)")
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

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
