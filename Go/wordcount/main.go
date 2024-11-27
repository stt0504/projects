//go:build !solution

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	stringCountMap := make(map[string]int)
	filenames := os.Args[1:]
	for _, filename := range filenames {
		data, err := os.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		stringsArr := strings.Split(string(data), "\n")
		for _, str := range stringsArr {
			stringCountMap[str]++
		}
	}
	for key, value := range stringCountMap {
		if value >= 2 {
			fmt.Printf("%d\t%s\n", value, key)
		}
	}
}
