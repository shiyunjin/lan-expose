package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"lan-expose/pkg/config"
)

func main() {
	configFile := flag.String("c", "upgrade.ini", "config file")
	flag.Parse()

	common, err := config.ParseUpgrade(*configFile)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
		return
	}

	listenAddress := common.Address + ":" + common.Port
	httpServer := &http.Server{
		Addr: listenAddress,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/favicon.ico", Favicon)
	mux.HandleFunc("/", Worker)

	httpServer.Handler = mux

	if common.SSL {
		go httpServer.ListenAndServeTLS(common.SSLCrt, common.SSLKey)
	} else {
		go httpServer.ListenAndServe()
	}

	log.Printf("Start LanExpose Upgrade Server on: %s\n", listenAddress)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	log.Println("shutdown server")
}
