package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	v1 "github.com/blabtm/v2k/platform/sm/rest/v1"
	"github.com/rs/cors"

	_ "github.com/blabtm/v2k/platform/sm/docker"
	_ "github.com/blabtm/v2k/platform/sm/swarm"
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

	mux.HandleFunc("GET /v1/svc", v1.Ls)
	mux.HandleFunc("GET /v1/svc/{name}/ps", v1.Ps)
	mux.HandleFunc("PUT /v1/svc/{name}/up", v1.Up)
	mux.HandleFunc("PUT /v1/svc/{name}/down", v1.Down)
	mux.HandleFunc("GET /v1/svc/{name}/config", v1.GetConfig)
	mux.HandleFunc("PUT /v1/svc/{name}/config", v1.PutConfig)

	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(mux)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Panic(err)
	}
}
