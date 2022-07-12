package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/shiyunjin/lan-expose/pkg/config"
	"github.com/shiyunjin/lan-expose/pkg/version"
)

func main() {
	configFile := flag.String("c", "proxy.ini", "config file")
	v := flag.Bool("v", false, "show version")
	flag.Parse()

	if *v {
		version.FormatFullVersion("Lan Expose Proxy")
		return
	}

	common, err := config.ParseProxy(*configFile)
	if err != nil {
		log.Fatalf("Fail to parse config file: %v", err)
		return
	}

	listenAddress := common.Address + ":" + common.Port

	log.Printf("Start Proxy Server on: %s\n", listenAddress)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		if err := ListenAndServe(listenAddress, common.SSLCrt, common.SSLKey, Worker); err != nil {
			log.Printf("[Serve][Error] listen and serve: %v\n", err)
			c <- os.Interrupt
		}
	}()

	<-c
	log.Println("Shutdown Server")
}
