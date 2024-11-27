//go:build !solution

package main

import (
	"flag"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var symbolMap = map[byte]string{
	'0': Zero,
	'1': One,
	'2': Two,
	'3': Three,
	'4': Four,
	'5': Five,
	'6': Six,
	'7': Seven,
	'8': Eight,
	'9': Nine,
}

func timeFormatIsValid(time string) bool {
	if len(time) != 8 || time[2] != ':' || time[5] != ':' {
		return false
	}

	parts := strings.Split(time, ":")
	if len(parts) != 3 {
		return false
	}

	for i, part := range parts {
		if len(part) != 2 || !unicode.IsDigit(rune(part[0])) || !unicode.IsDigit(rune(part[1])) {
			return false
		}

		if i == 0 {
			if part[0] < '0' || part[0] > '2' || (part[0] == '2' && part[1] > '3') {
				return false
			}
		}

		if i == 1 || i == 2 {
			if part[0] < '0' || part[0] > '5' || (part[0] == '5' && part[1] > '9') {
				return false
			}
		}
	}

	return true
}

func main() {
	port := flag.String("port", "6029", "port to listen on")
	flag.Parse()

	http.HandleFunc("/", handler)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	timeParam := r.URL.Query().Get("time")
	kParam := r.URL.Query().Get("k")

	if timeParam == "" {
		timeParam = time.Now().Format("00:00:00")
	} else if !timeFormatIsValid(timeParam) {
		http.Error(w, "invalid time format", http.StatusBadRequest)
		return
	}

	k := 1
	if kParam != "" {
		var err error
		k, err = strconv.Atoi(kParam)
		if err != nil || k < 1 || k > 30 {
			http.Error(w, "invalid k", http.StatusBadRequest)
			return
		}
	}

	img := generateClockImage(timeParam, k)

	w.Header().Set("Content-Type", "image/png")
	err := png.Encode(w, img)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func generateClockImage(time string, k int) image.Image {
	parts := strings.Split(time, ":")
	width := 6*8 + 2*4
	height := 12

	img := image.NewRGBA(image.Rect(0, 0, width*k, height*k))
	begin := 0
	for i, part := range parts {
		if i > 0 {
			drawColon(img, k, begin)
			begin += 4 * k
		}
		for j := 0; j < 2; j++ {
			drawDigit(img, symbolMap[part[j]], k, begin)
			begin += 8 * k
		}
	}

	return img
}

func drawColon(img *image.RGBA, k int, begin int) {
	height := 12
	width := 4
	symbol := Colon
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			for dx := 0; dx < k; dx++ {
				for dy := 0; dy < k; dy++ {
					if string(symbol[row*5+col]) == "1" {
						img.Set(begin+col*k+dx, row*k+dy, Cyan)
					} else {
						img.Set(begin+col*k+dx, row*k+dy, color.White)
					}
				}
			}
		}
	}
}

func drawDigit(img *image.RGBA, s string, k int, begin int) {
	height := 12
	width := 8
	symbol := s
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			for dx := 0; dx < k; dx++ {
				for dy := 0; dy < k; dy++ {
					if string(symbol[row*9+col]) == "1" {
						img.Set(begin+col*k+dx, row*k+dy, Cyan)
					} else {
						img.Set(begin+col*k+dx, row*k+dy, color.White)
					}
				}
			}
		}
	}
}
