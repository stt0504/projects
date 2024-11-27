//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	urls := os.Args[1:]

	if len(urls) == 0 {
		fmt.Printf("%s", "No urls were provided")
		return
	}

	var wg sync.WaitGroup
	num := 0

	startAll := time.Now()
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			start := time.Now()
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				return
			}
			body, _ := io.ReadAll(resp.Body)
			num++
			elapsed := time.Since(start).Seconds()
			fmt.Printf("%d. %fs\t%d\t%s\n", num, elapsed, len(body), url)
			resp.Body.Close()
		}(url)
	}

	wg.Wait()
	elapsedAll := time.Since(startAll).Seconds()
	fmt.Printf("   %fs elapsed\n", elapsedAll)
}
