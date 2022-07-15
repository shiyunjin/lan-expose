package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/netinternet/remoteaddr"
	"github.com/shiyunjin/lan-expose/pkg/config"
)

//go:embed check-page/index.html
//go:embed check-page/favicon.ico
var f embed.FS

var (
	checkT     *template.Template
	checkTHave = false
)

func init() {
	var err error
	checkT, err = template.ParseFS(f, "check-page/index.html")
	if err != nil {
		log.Printf("internal check template error: %v", err)
		return
	}
	checkTHave = true
}

func Worker(w http.ResponseWriter, r *http.Request) {
	// ext info
	ip, _ := remoteaddr.Parse().IP(r)
	requestId := r.Header.Get("CF-ray")
	if requestId == "" {
		requestId = GetRay()
	}

	// search service
	service, ok := config.SearchServiceWithDomain(r.Host)
	if !ok {
		w.WriteHeader(404)
		_, _ = w.Write([]byte("domain not found"))
		log.Printf("[Request][Miss] %s : %s[%s]", r.Host, ip, requestId)
		return
	}

	// upgrade h3
	altAddr := service.DestDomain + ":" + service.DestPort
	w.Header().Add("Alt-Svc",
		`h3="`+altAddr+`"; ma=2592000; persist=1,`+
			`h3-29="`+altAddr+`"; ma=2592000; persist=1`,
	)

	// upgrade websocket proxy
	if WorkerWebSocket(w, r, service, ip, requestId) {
		return
	}

	// send default html
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(200)

	log.Printf("[Request][Hit] %s => %s : %s[%s]", r.Host, altAddr, ip, requestId)

	// send html
	if !checkTHave {
		_, _ = w.Write([]byte("This service requires use of the HTTP/3.0 protocol, wait reload.<script>location.reload();</script>"))
	}

	checkMS := config.GetServerUpgradeCommon().CheckMS
	if service.CheckMS != 0 {
		checkMS = service.CheckMS
	}

	if err := checkT.Execute(w, struct {
		IP        string
		RequestID string
		CheckMS   uint32
	}{
		IP:        ip,
		RequestID: requestId,
		CheckMS:   checkMS,
	}); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

// Favicon icon file
func Favicon(writer http.ResponseWriter, request *http.Request) {
	content, err := f.ReadFile("check-page/favicon.ico")
	if err != nil {
		writer.WriteHeader(404)
		_, _ = writer.Write([]byte("file not found"))
		return
	}
	writer.WriteHeader(200)
	_, _ = writer.Write(content)
}
