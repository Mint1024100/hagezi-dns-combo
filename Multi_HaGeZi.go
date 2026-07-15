package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	urls := []string{
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/pro.plus-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/tif-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/spam-tlds-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/anti.piracy-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/native.amazon-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/native.apple-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/native.winoffice-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/native.samsung-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/native.huawei-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/gambling-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/nsfw-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/fake-onlydomains.txt",
		"https://raw.githubusercontent.com/hagezi/dns-blocklists/main/wildcard/popupads-onlydomains.txt",
	}

	ch := make(chan string, 100000)
	var fileWg sync.WaitGroup

	fileWg.Add(1)
	go func() {
		defer fileWg.Done()
		seen := make(map[string]struct{})
		file, err := os.Create("M_HaGeZi.txt")
		if err != nil {
			fmt.Println("Error creating file:", err)
			for range ch {
			}
			return
		}
		defer file.Close()

		for url := range ch {
			createURL := strings.TrimSpace(url)
			if _, exists := seen[createURL]; !exists {
				seen[createURL] = struct{}{}
				if _, err := file.WriteString(createURL + "\n"); err != nil {
					fmt.Println("Error writing to file:", err)
				}
			}
		}
	}()

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			fetchURL(u, ch)
		}(url)
	}
	wg.Wait()
	close(ch)
	fileWg.Wait()
}

func fetchURL(url string, ch chan<- string) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		ch <- line
	}
}
