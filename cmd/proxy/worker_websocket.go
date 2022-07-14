package main

import (
	"log"
	"net/http"

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

	Worker(w, r)
}
