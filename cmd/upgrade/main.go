package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/shiyunjin/lan-expose/pkg/config"
	"github.com/shiyunjin/lan-expose/pkg/version"
)

func main() {
	configFile := flag.String("c", "upgrade.ini", "config file")
	v := flag.Bool("v", false, "show version")
	flag.Parse()

	if *v {
		fmt.Println(version.FormatFullVersion("Lan Expose Upgrade"))
		return
	}

	common, err := config.ParseUpgrade(*configFile)
	if err != nil {
		log.Fatalf("Fail to parse config file: %v", err)
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
	log.Println("Shutdown Server")
}
