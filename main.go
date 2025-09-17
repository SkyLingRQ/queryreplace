package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
)

func queryReplace(file string, text string, filesave string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()

	content, err := os.Open(file)
	if err != nil {
		fmt.Println("\033[31m[x] Error reading file.\033[0m")
		return
	}

	urls := []string{}

	defer content.Close()

	bufioScanner := bufio.NewScanner(content)
	for bufioScanner.Scan() {
		URL := bufioScanner.Text()

		if !strings.Contains(URL, "=") {
			continue
		}

		urlParsed, err := url.Parse(URL)
		if err != nil {
			continue
		}

		qs := urlParsed.Query()
		for key := range qs {
			qs.Set(key, text)
		}
		urlParsed.RawQuery = qs.Encode()
		urls = append(urls, urlParsed.String())
	}

	fileCreate, err := os.Create(filesave)
	if err != nil {
		fmt.Println("\033[31m[x] Error creating file.\033[0m")
		return
	}

	defer fileCreate.Close()

	for _, url := range urls {
		fmt.Fprintln(fileCreate, url)
	}

}

func main() {
	var wg sync.WaitGroup
	file := flag.String("file", "urls.txt", "File with URLs to scan.")
	textReplace := flag.String("text", "<script>alert(8);</script>", "Text for replace the query values.")
	fileSave := flag.String("o", "urls_queryreplace.txt", "Name of file save.")

	flag.Parse()

	sem := make(chan struct{}, 50)

	wg.Add(1)
	go queryReplace(*file, *textReplace, *fileSave, &wg, sem)
	wg.Wait()
}
