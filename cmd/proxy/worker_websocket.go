package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/shiyunjin/lan-expose/pkg/config"
	"github.com/shiyunjin/lan-expose/pkg/utils"
)

func WorkerWebSocket(w http.ResponseWriter, r *http.Request) {
	reqUpType := utils.WebSocketUpgradeType(r.Header)

	if reqUpType == "" {
		log.Printf("[Block] Disallow direct access to the proxy service : %s %s", r.Host, r.URL)
		w.WriteHeader(403)
		w.Write([]byte("access block"))
		return
	}

	webSocketConfig := config.GetServerProxyWebSocket()
	if webSocketConfig.Mode302Domain != "" && strings.HasPrefix(r.URL.String(), config.ServerWebSocketMode302PrefixURI) {
		uriArray := strings.Split(r.URL.String(), "/")
		r.Host = uriArray[2]
		if len(uriArray) > 3 {
			uriArray = uriArray[3:len(uriArray)]
		} else {
			uriArray = []string{}
		}

		u, err := url.Parse("/" + strings.Join(uriArray, "/"))
		if err != nil {
			log.Printf("[Error] %s:%s %v", r.Host, r.URL, err)
			w.WriteHeader(502)
			_, _ = w.Write([]byte("internal error - code: pwwup"))
			return
		}

		r.URL = u
	}

	Worker(w, r)
}
