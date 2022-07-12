package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/shiyunjin/lan-expose/pkg/config"
)

func Worker(w http.ResponseWriter, r *http.Request) {
	u, ok := config.SearchTargetWithDomain(r.Host)
	if !ok {
		w.WriteHeader(404)
		_, _ = w.Write([]byte("target not found"))
		log.Printf("[Request][Miss] %s : %s\n", r.Host, r.URL)
		return
	}

	log.Printf("[Request][Hit] %s => %s : %s\n", r.Host, u.String(), r.URL)

	reverseProxy := httputil.NewSingleHostReverseProxy(u)

	reverseProxy.ServeHTTP(w, r)
}
