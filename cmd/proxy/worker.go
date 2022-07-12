package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Worker(w http.ResponseWriter, r *http.Request) {
	targetUrl, ok := hostMap[r.Host]
	if !ok {
		w.WriteHeader(404)
		log.Printf("[Request][Miss] %s : %s\n", r.Host, r.URL)
		w.Write([]byte("404 not found"))
		return
	}

	u, err := url.Parse(targetUrl)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte("internal error"))
		log.Printf("[Request][Error] %s => %s : %s : url parse error\n", r.Host, targetUrl, r.URL)
		return
	}
	log.Printf("[Request][Hit] %s => %s : %s\n", r.Host, targetUrl, r.URL)

	reverseProxy := httputil.NewSingleHostReverseProxy(u)

	reverseProxy.ServeHTTP(w, r)
}
