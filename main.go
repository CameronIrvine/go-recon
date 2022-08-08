package main

import (
	"fmt"
	"regexp"
	"strings"

	term "github.com/buger/goterm"
)

func main() {
	for {
		term.Clear()
		term.MoveCursor(1, 1)
		var host string
		fmt.Print("Enter a domain to search (or 'q' to exit): ")
		fmt.Scanln(&host)

		r, _ := regexp.Compile("^(?:\\*\\.)?[a-z0-9]+(?:[\\-.][a-z0-9]+)*\\.[a-z]{2,6}$")
		if host == "q" {
			term.Flush()
			term.MoveCursor(1, 1)
			break
		} else if !r.MatchString(host) {
			continue
		}
		term.Flush()
		term.MoveCursor(1, 1)
		subdomains := FetchSubdomains(&host)
		for i, record := range subdomains {
			fmt.Printf("%3d. %v (%v)\n", i+1, record.subdomain, strings.Join(record.ipAddresses, ", "))
		}
		fmt.Println()
	}
}
