package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	v1 "github.com/blabtm/v2k/platform/sm/rest/v1"
	"github.com/rs/cors"

	_ "github.com/blabtm/v2k/platform/sm/compose"
	_ "github.com/blabtm/v2k/platform/sm/swarm"
	_ "github.com/blabtm/v2k/platform/sm/systemd"
)

var addr string
var allowedOrigins []string

func init() {
	addr = os.Getenv("ADDR")

	if addr == "" {
		addr = ":80"
	}

	allowedOrigins = strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
}

func main() {
	mux := http.NewServeMux()

	// List all registered services.
	mux.HandleFunc("GET /api", v1.Ls)

	// Get the current operational status of the service.
	mux.HandleFunc("GET /api/{name}/ps", v1.Ps)

	// Get the runtime (desired) configuration of the service.
	mux.HandleFunc("GET /api/{name}", v1.Get)

	// Set the runtime (desired) configuration of the service.
	mux.HandleFunc("PUT /api/{name}", v1.Set)

	// Fork oneshot service, merging and setting desired configuration.
	mux.HandleFunc("PUT /api/{name}/fork", v1.Fork)

	mux.HandleFunc("GET /api/{name}/ui", v1.GetUI)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(mux)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Panic(err)
	}
}
