package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// Get env var or default
func getEnvOptional(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	fmt.Printf("No value specified for %s\n", key)
	log.Fatal()
	return ""
}

func main() {

	portSuffix := ":" + getEnvOptional("PORT", "8080")
	registryUsername := getEnv("USERNAME")
	registryPassword := getEnv("PASSWORD")
	registry := getEnv("REGISTRY")

	target, err := url.Parse("https://" + registry + ":443")
	log.Printf("forwarding to -> %s://%s\n", target.Scheme, target.Host)

	if err != nil {
		log.Fatal(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("req.URL.Host=%s req.RequestURI=%s\n", req.URL.Host, req.RequestURI)
		req.Host = registry
		if strings.HasPrefix(req.RequestURI, "/v2/") == true {
			req.SetBasicAuth(registryUsername, registryPassword)
		}

		proxy.ServeHTTP(w, req)
	})

	err = http.ListenAndServe(portSuffix, nil)
	if err != nil {
		panic(err)
	}

	log.Printf("Listening on" + portSuffix)
}
