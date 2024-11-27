//go:build !solution

package main

import (
	"bytes"
	"flag"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type Rule struct {
	Endpoint               string   `yaml:"endpoint"`
	ForbiddenUserAgents    []string `yaml:"forbidden_user_agents"`
	ForbiddenHeaders       []string `yaml:"forbidden_headers"`
	RequiredHeaders        []string `yaml:"required_headers"`
	MaxRequestLengthBytes  int      `yaml:"max_request_length_bytes"`
	MaxResponseLengthBytes int      `yaml:"max_response_length_bytes"`
	ForbiddenResponseCodes []int    `yaml:"forbidden_response_codes"`
	ForbiddenRequestRe     []string `yaml:"forbidden_request_re"`
	ForbiddenResponseRe    []string `yaml:"forbidden_response_re"`
}

type Config struct {
	Rules []Rule `yaml:"rules"`
}

func parseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func matchesAny(patterns []string, value string) bool {
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, value); matched {
			return true
		}
	}
	return false
}

func shouldBlockRequest(r *http.Request, rule Rule) bool {
	if rule.MaxRequestLengthBytes > 0 && r.ContentLength > int64(rule.MaxRequestLengthBytes) {
		return true
	}

	for _, header := range rule.RequiredHeaders {
		if r.Header.Get(header) == "" {
			return true
		}
	}

	for _, forbiddenHeader := range rule.ForbiddenHeaders {
		for name, values := range r.Header {
			for _, value := range values {
				if forbiddenHeader == name+": "+value {
					return true
				}
			}
		}
	}

	if matchesAny(rule.ForbiddenUserAgents, r.UserAgent()) {
		return true
	}
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return matchesAny(rule.ForbiddenRequestRe, string(body))
}

func shouldBlockResponse(resp *http.Response, rule Rule) bool {
	if rule.MaxResponseLengthBytes > 0 && resp.ContentLength > int64(rule.MaxResponseLengthBytes) {
		return true
	}
	for _, code := range rule.ForbiddenResponseCodes {
		if resp.StatusCode == code {
			return true
		}
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	return matchesAny(rule.ForbiddenResponseRe, string(body))
}

func main() {
	serviceAddr := flag.String("service-addr", "", "service address")
	confPath := flag.String("conf", "", "config path")
	addr := flag.String("addr", ":8081", "firewall address")
	flag.Parse()

	config, err := parseConfig(*confPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	URL, err := url.Parse(*serviceAddr)
	if err != nil {
		log.Fatalf("%v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(URL)
	proxy.ModifyResponse = func(resp *http.Response) error {
		for _, rule := range config.Rules {
			if strings.HasPrefix(resp.Request.URL.Path, rule.Endpoint) && shouldBlockResponse(resp, rule) {
				resp.Body = io.NopCloser(strings.NewReader("Forbidden"))
				resp.StatusCode = http.StatusForbidden
				resp.ContentLength = 9
				resp.Header.Set("Content-Length", "9")
				return nil
			}
		}
		return nil
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, rule := range config.Rules {
			if strings.HasPrefix(r.URL.Path, rule.Endpoint) && shouldBlockRequest(r, rule) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
		proxy.ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
