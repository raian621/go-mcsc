package main

import (
	"log"

	"github.com/raian621/minecraft-server-controller/minecraft"
	"github.com/raian621/minecraft-server-controller/webserver"
)

func main() {
	log.Println("creating server folder...")
	if err := minecraft.CreateServerFolder("server-data"); err != nil {
		log.Fatalln(err)
	}

	log.Println("loading server config...")
	if err := minecraft.LoadConfig("server-data/server-config.json"); err != nil {
		log.Fatalln(err)
	}

	webserver.StartWebServer("localhost", "5000", "", "")
}
