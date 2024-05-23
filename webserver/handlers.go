package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/raian621/minecraft-server-controller/minecraft"
)

func addHandlers(server *http.Server) {
	mux := http.NewServeMux()

	mux.HandleFunc("/start", POST(handleServerStart))
	mux.HandleFunc("/stop", POST(handleServerStop))
	mux.HandleFunc("/command", POST(handleServerCommand))

	server.Handler = mux
}

// POST
func handleServerStart(w http.ResponseWriter, r *http.Request) {
	log.Println("starting minecraft server...")

	if err := minecraft.StartMinecraftServer(); err != nil {
		log.Println(err)

		errMsg := MessageResponse{
			Message: "could not start minecraft server",
		}

		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(errMsg)
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

// POST
func handleServerStop(w http.ResponseWriter, r *http.Request) {
	log.Println("stopping minecraft server")

	minecraft.StopMinecraftServer()
}

// POST
func handleServerCommand(w http.ResponseWriter, r *http.Request) {
	if minecraft.ServerConsole == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	var command struct {
		Cmd string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var err error
	maxTries := 3

	for tries := 0; tries < maxTries+1; tries++ {
		if locked := minecraft.ServerConsole.TryLock(); locked {
			_, err = fmt.Fprintln(minecraft.ServerConsole.Stdin, command.Cmd)
			minecraft.ServerConsole.Unlock()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			break
		} else {
			if tries == maxTries {
				w.WriteHeader(http.StatusRequestTimeout)
			}

			time.Sleep(time.Duration(math.Pow(10, float64(tries)/2)) * time.Second)
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
