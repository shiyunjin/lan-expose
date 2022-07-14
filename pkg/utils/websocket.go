package utils

import (
	"net/http"

	"golang.org/x/net/http/httpguts"
)

func WebSocketUpgradeType(h http.Header) string {
	if !httpguts.HeaderValuesContainsToken(h["Connection"], "Upgrade") {
		return ""
	}
	return h.Get("Upgrade")
}
