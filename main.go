package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/raian621/go-mcsc/api"
	"github.com/raian621/go-mcsc/minecraft"
)

func main() {
	host := flag.String("host", "0.0.0.0", "host to serve requests from")
	port := flag.String("port", "5000", "port to listen on")
	flag.Parse()

	log.Println("creating server folder...")
	if err := minecraft.CreateServerFolder("server-data"); err != nil {
		log.Fatalln(err)
	}

	mcServer := minecraft.NewJavaMinecraftServer(&minecraft.MinecraftServerConfigFilepaths{
		Allowlist:          "server-data/whitelist.json",
		Args:               "server-data/args.json",
		BannedIPs:          "server-data/banned-ips.json",
		BannedPlayers:      "server-data/banned-players.json",
		Config:             "server-data/config.json",
		Ops:                "server-data/ops.json",
		Properties:         "server-data/properties.json",
		PropertiesTemplate: "templates/server.properties.tmpl",
		Versions:           "data/server-download-links.json",
	})

	log.Println("loading server config...")
	if err := mcServer.LoadConfigs(); err != nil {
		log.Fatalln(err)
	}

	server := api.NewServerController(mcServer)
	r := chi.NewMux()
	h := api.HandlerFromMux(server, r)

	addr := net.JoinHostPort(*host, *port)
	s := &http.Server{
		Handler: h,
		Addr:    addr,
	}

	log.Printf("hosting server controller at http://%s\n", addr)
	log.Fatal(s.ListenAndServe())
}
