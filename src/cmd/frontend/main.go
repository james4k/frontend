package main

import (
	_ "expvar"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)
	go func() {
		// immediately close all https connections..works better for
		// most OSes it seems
		l, err := net.Listen("tcp", ":8443")
		if err != nil {
			return
		}
		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			conn.Close()
		}
	}()
	server := &http.Server{
		Addr:         ":8000",
		Handler:      nil,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	http.Handle("/", assemble())
	log.Fatal(server.ListenAndServe())
}
