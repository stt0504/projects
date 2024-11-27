package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
)

var (
	mu     sync.Mutex
	urlMap = make(map[string]string)
)

func generateKey(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func main() {
	port := flag.String("port", "6029", "port to listen on")
	flag.Parse()

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/go/", goHandler)

	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]string

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	key := generateKey(req["url"])
	urlMap[key] = req["url"]

	resp := map[string]string{
		"url": req["url"],
		"key": key,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

func goHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[4:]

	mu.Lock()
	defer mu.Unlock()

	url, found := urlMap[key]
	if !found {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
