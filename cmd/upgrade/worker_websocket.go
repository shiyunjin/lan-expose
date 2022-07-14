package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/shiyunjin/lan-expose/pkg/utils"
)

func WorkerWebSocket(w http.ResponseWriter, r *http.Request, domain, port, ip, requestId string) bool {
	reqUpType := utils.WebSocketUpgradeType(r.Header)

	if reqUpType == "" {
		return false
	}

	log.Printf("[Request][Proxy] %s => %s:%s : %s[%s] %s", r.Host, domain, port, ip, requestId, r.URL)

	u, err := url.Parse("https://" + domain + ":" + port + "/")
	if err != nil {
		log.Printf("[Error] %s:%s %v", domain, port, err)
		w.WriteHeader(502)
		_, _ = w.Write([]byte("internal error - code: uwwup"))
		return true
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	reverseProxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	reverseProxy.ServeHTTP(w, r)

	return true
}
