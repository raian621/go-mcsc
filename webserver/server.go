package webserver

import (
	"log"
	"net"
	"net/http"
)

func StartWebServer(addr string, port string, keyfile, certfile string) {
	server := &http.Server{
		Addr: net.JoinHostPort(addr, port),
	}
	addHandlers(server)

	log.Fatalln(server.ListenAndServe())
}
