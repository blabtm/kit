package main

import (
	"log"
	"net/http"

	v1 "github.com/blabtm/v2k/platform/sm/rest/v1"

	_ "github.com/blabtm/v2k/platform/sm/docker"
	_ "github.com/blabtm/v2k/platform/sm/swarm"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/svc", v1.Ls)
	mux.HandleFunc("GET /v1/svc/{name}/ps", v1.Ps)
	mux.HandleFunc("PUT /v1/svc/{name}/up", v1.Up)
	mux.HandleFunc("PUT /v1/svc/{name}/down", v1.Down)
	mux.HandleFunc("GET /v1/svc/{name}/config", v1.GetConfig)
	mux.HandleFunc("PUT /v1/svc/{name}/config", v1.PutConfig)

	if err := http.ListenAndServe(":80", mux); err != nil {
		log.Panic(err)
	}
}
