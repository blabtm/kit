package main

import (
	"log"
	"net/http"
	"os"

	v1 "github.com/blabtm/v2k/platform/sm/rest/v1"

	_ "github.com/blabtm/v2k/platform/sm/docker"
)

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(os.Getenv("DOCS_PATH")))

	mux.Handle("/v1/docs/", http.StripPrefix("/v1/docs/", fs))

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
