//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	urls := os.Args[1:]
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		body, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			panic(err2)
		}
		fmt.Printf("%s", body)

		err = resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}
}
