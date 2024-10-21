package internal

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"

	"github.com/devdouglasferreira/stockscrawler/internal/models"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

func FetchURL(url string) (*http.Response, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}
	req, _ := http.NewRequest("GET", url, nil)
	setHeaders(req)
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("failed to fetch %s: %s", url, resp.Status)
	}

	return resp, nil
}

func ParseHTML(httpResponse *http.Response) (*models.StockPrice, error) {

	file, err := os.Create("body.html")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var bodyBuffer []byte
	httpResponse.Body.Read(bodyBuffer)
	file.Write(bodyBuffer)

	defer file.Close()
	defer httpResponse.Body.Close()

	doc, err := html.Parse(httpResponse.Body)

	if err != nil {
		log.Fatal(err)
	}

	var openPrice string
	var closePrice string
	var highPrice string
	var lowPrice string
	var volume string

	var fNode func(*html.Node)

	fInnerNode := func(n *html.Node, tag string, a *string) {

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == tag {
				extracted := extractInnerText(c)
				*a = extracted
			}
			fNode(c)
		}
	}

	var fInnerTable func(n *html.Node)
	fInnerTable = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "td" {
				fInnerTable(c)
			}
			if c.Type == html.TextNode && c.Data == "Abertura" {
				openPrice = extractInnerText(c.Parent.NextSibling.FirstChild)
			}
			if c.Type == html.TextNode && c.Data == "Volume" {
				volume = extractInnerText(c.Parent.NextSibling.FirstChild)
			}
			fInnerTable(c)
		}
	}
	fNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {

			for _, attr := range n.Attr {

				if attr.Key == "class" && strings.Contains(attr.Val, "minimo") {
					fInnerNode(n, "p", &lowPrice)
				}

				if attr.Key == "class" && strings.Contains(attr.Val, "maximo") {
					fInnerNode(n, "p", &highPrice)
				}

				if attr.Key == "class" && strings.Contains(attr.Val, "value") {
					fInnerNode(n, "p", &closePrice)
				}

				if attr.Key == "class" && strings.Contains(attr.Val, "tables") {
					fInnerTable(n.FirstChild.NextSibling)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fNode(c)
		}
	}
	fNode(doc.FirstChild.NextSibling.LastChild)

	o, _ := strconv.ParseFloat(strings.Replace(strings.Replace(openPrice, ",", ".", -1), " ", "", -1), 64)
	c, _ := strconv.ParseFloat(strings.Replace(closePrice, ",", ".", -1), 64)
	h, _ := strconv.ParseFloat(strings.Replace(highPrice, ",", ".", -1), 64)
	l, _ := strconv.ParseFloat(strings.Replace(lowPrice, ",", ".", -1), 64)

	stockPrice := models.StockPrice{Ticker: "", Open: o, Close: c, High: h, Low: l, Volume: convertVolumeStr(volume)}

	defer httpResponse.Body.Close()
	return &stockPrice, nil
}

func convertVolumeStr(strVol string) int64 {

	var multiplier float64

	if strings.Contains(strVol, "M") {
		multiplier = 1000000
	} else if strings.Contains(strVol, "B") {
		multiplier = 1000000000
	}

	volStr := strings.Replace(strVol, ",", ".", -1)
	volStr = strings.Replace(volStr, "$ ", "", -1)
	volStr = strings.Replace(volStr, " M", "", -1)
	volStr = strings.Replace(volStr, " B", "", -1)
	v, _ := strconv.ParseFloat(volStr, 64)

	return int64(v * multiplier)
}

func extractInnerText(n *html.Node) string {
	var result string
	var f func(*html.Node)

	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			result += n.Data
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(n)
	return result
}

func setHeaders(req *http.Request) {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "PostmanRuntime/7.41.1")
}
