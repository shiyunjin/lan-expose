package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var hostMap = map[string]string{
	"uptime.test.cn": "http://uptime.test.cn/",
}

func main() {
	port := flag.String("port", "690", "listen port")
	sslFile := flag.String("ssl-file", "./ssl.crt", "ssl cert file")
	sslKey := flag.String("ssl-key", "./ssl.key", "ssl cert key")
	flag.Parse()

	log.Printf("start proxy server on: %s\n", *port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		if err := ListenAndServe(*port, *sslFile, *sslKey, Worker); err != nil {
			log.Printf("[Serve][Error] listen and serve: %v\n", err)
			c <- os.Interrupt
		}
	}()

	<-c
	log.Println("shutdown server")
}
