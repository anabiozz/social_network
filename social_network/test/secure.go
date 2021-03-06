package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

const (
	serviceName  = "localhost"
	SSLport      = ":443"
	HTTPport     = ":8080"
	SSLProtocol  = "https://"
	HTTPprotocol = "http://"
)

func secureRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You have arrived at port 443, but you are not yet secure")
}

func redirectNonSecure(w http.ResponseWriter, r *http.Request) {
	log.Println("Non-secure request initiated, redirection")
	redirectURL := SSLProtocol + serviceName + r.RequestURI
	http.Redirect(w, r, redirectURL, http.StatusOK)
}

func main() {
	wg := sync.WaitGroup{}
	log.Println("Starting redirection, server try to access @ http:")

	wg.Add(1)
	go func() {
		http.ListenAndServe(HTTPport, http.HandlerFunc(redirectNonSecure))
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		http.ListenAndServe(SSLport, http.HandlerFunc(secureRequest))
		wg.Done()
	}()
	wg.Wait()
}
