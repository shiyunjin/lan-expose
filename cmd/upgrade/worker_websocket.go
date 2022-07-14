package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/shiyunjin/lan-expose/pkg/config"
	"github.com/shiyunjin/lan-expose/pkg/utils"
)

func WorkerWebSocket(w http.ResponseWriter, r *http.Request, service config.Service, ip, requestId string) bool {
	reqUpType := utils.WebSocketUpgradeType(r.Header)

	if reqUpType == "" {
		return false
	}

	switch service.WebSocketMode {
	case config.ServerWebSocketModeProxy:
		log.Printf("[Request][Proxy] %s => %s:%s : %s[%s] %s", r.Host, service.DestDomain, service.DestPort, ip, requestId, r.URL)

		u, err := url.Parse("https://" + service.DestDomain + ":" + service.DestPort + "/")
		if err != nil {
			log.Printf("[Error] %s:%s %v", service.DestDomain, service.DestPort, err)
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
	case config.ServerWebSocketMode302:
		log.Printf("[Request][Redirect] %s => %s:%s : %s[%s] %s", r.Host, service.WebSocketMode302Domain, service.DestPort, ip, requestId, r.URL)

		w.Header().Add("Location", "wss://"+service.WebSocketMode302Domain+":"+service.DestPort+config.ServerWebSocketMode302PrefixURI+r.Host+r.URL.String())
		w.WriteHeader(302)

		return true
	case config.ServerWebSocketModeBlock:
		log.Printf("[Request][Block] %s => %s:%s : %s[%s] %s", r.Host, service.DestDomain, service.DestPort, ip, requestId, r.URL)
	default:
		log.Printf("[Request][Block][undefined-mode] %s => %s:%s : %s[%s] %s", r.Host, service.DestDomain, service.DestPort, ip, requestId, r.URL)
	}

	w.WriteHeader(403)
	w.Write([]byte("block websocket mode"))
	return true
}
