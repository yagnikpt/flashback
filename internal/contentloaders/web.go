package contentloaders

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	h2m "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
)

func formatHead(head *goquery.Selection) string {
	var lines []string
	head.Children().Each(func(i int, s *goquery.Selection) {
		tag := s.Get(0).Data
		var attrs []string
		for _, attr := range s.Get(0).Attr {
			attrs = append(attrs, fmt.Sprintf("%s=%s", attr.Key, attr.Val))
		}
		if tag == "title" {
			text := s.Text()
			attrs = append(attrs, fmt.Sprintf("text=%s", text))
		}
		if len(attrs) > 0 {
			line := fmt.Sprintf("[%s] %s", tag, strings.Join(attrs, " "))
			lines = append(lines, line)
		}
	})
	return strings.Join(lines, "\n\n")
}

func GetWebPage(target string) (string, error) {
	target = strings.TrimRight(target, ".,;:!?)")
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "https://" + target
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "facebookexternalhit/1.1")
	req.Header.Set("sec-ch-ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)

	for {
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location := resp.Header.Get("Location")
			resp.Body.Close()
			if location == "" {
				return "", fmt.Errorf("redirect without location")
			}

			// Resolve relative URLs
			if !strings.HasPrefix(location, "http") {
				baseURL, err := url.Parse(target)
				if err != nil {
					return "", err
				}
				resolved := baseURL.ResolveReference(&url.URL{Path: location})
				location = resolved.String()
			}

			target = location
		} else if resp.StatusCode == 200 {
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return "", err
			}

			doc.Find("script").Remove()
			doc.Find("style").Remove()
			doc.Find("link[rel='stylesheet']").Remove()
			doc.Find("noscript").Remove()
			doc.Find("svg").Remove()
			doc.Find("nav").Remove()
			doc.Find("header").Remove()
			doc.Find("footer").Remove()
			doc.Find("aside").Remove()
			doc.Find("dialog").Remove()
			doc.Find("*[role='dialog']").Remove()

			// Remove inline style attributes from all elements
			doc.Find("*").Each(func(i int, s *goquery.Selection) {
				s.RemoveAttr("style")
			})

			// Remove elements with 'error' in any attribute value
			doc.Find("*").Each(func(i int, s *goquery.Selection) {
				node := s.Get(0)
				for _, attr := range node.Attr {
					if strings.Contains(strings.ToLower(attr.Val), "error") || strings.Contains(attr.Val, "blob") {
						s.Remove()
						return
					}
				}
			})

			head := doc.Find("head")
			formattedHead := formatHead(head)

			body := doc.Find("body")
			bodyHtml, err := body.Html()
			if err != nil {
				return "", err
			}

			bodyMarkdown, err := h2m.ConvertString(bodyHtml)
			if err != nil {
				return "", err
			}

			return formattedHead + "\n" + bodyMarkdown, nil
		} else {
			resp.Body.Close()
			return "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
		}
	}
}
