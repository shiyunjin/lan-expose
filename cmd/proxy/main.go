package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"lan-expose/pkg/config"
)

func main() {
	configFile := flag.String("c", "proxy.ini", "config file")
	flag.Parse()

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
