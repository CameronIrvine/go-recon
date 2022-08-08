package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	term "github.com/buger/goterm"
)

type subdomainRecord struct {
	subdomain   string
	ipAddresses []string
}

func FetchSubdomains(host *string) []subdomainRecord {
	fmt.Println(term.Color(fmt.Sprintf("Fetching subdomains for %v", *host), term.GREEN))
	html := getHtml(host)
	subdomains := parseHtml(host, html)
	records := []subdomainRecord{}
	for _, subdomain := range subdomains {
		newRecord := subdomainRecord{subdomain: subdomain}
		newRecord.ipAddresses = getIpAddresses(subdomain)
		records = append(records, newRecord)
	}
	return records
}

func getHtml(host *string) io.Reader {
	dnsDumpsterUrl := "https://dnsdumpster.com"

	// do an initial GET on the home page to get the CSRF token
	res, err := http.Get(dnsDumpsterUrl)
	if err != nil {
		log.Fatal("Failed to load DNSdumpster home page")
	}
	csrfToken := res.Cookies()[0]

	// build POST payload
	payload := url.Values{}
	payload.Set("csrfmiddlewaretoken", csrfToken.Value)
	payload.Set("targetip", *host)
	payload.Set("user", "free")

	// build and execute POST request
	client := http.Client{}
	req, _ := http.NewRequest("POST", dnsDumpsterUrl, strings.NewReader(payload.Encode()))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", dnsDumpsterUrl)
	req.AddCookie(csrfToken)
	res, err = client.Do(req)

	if err != nil {
		log.Fatal("Failed to load subdomains for the selected host")
	}
	return res.Body
}

func parseHtml(host *string, body io.Reader) []string {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal("Failed to create document")
	}
	subdomains := []string{}
	doc.Find("td.col-md-4").Each(func(i int, selection *goquery.Selection) {
		r, _ := regexp.Compile(".*." + *host)
		text := selection.Text()
		subdomain := r.FindString((text))
		if subdomain != "" {
			subdomains = append(subdomains, subdomain)
		}
	})
	return subdomains
}

func getIpAddresses(subdomain string) []string {
	out, err := exec.Command("dig", subdomain, "+short").Output()
	if err != nil {
		log.Fatal(err)
	}
	addresses := []string{}
	for _, address := range strings.Split(string(out), "\n") {
		strings.TrimSpace(address)
		if address != "" {
			addresses = append(addresses, address)
		}
	}
	return addresses
}
