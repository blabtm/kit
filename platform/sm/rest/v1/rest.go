package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/blabtm/v2k/platform/sm/service"
)

type RunRequest struct {
	Message string
	Diff    string
}

type RunResponse struct {
	IID string
}

func Ls(w http.ResponseWriter, r *http.Request) {
	services, err := service.Ls()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("ls: %v", err).Error()))
	}

	out, err := json.Marshal(services)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("json: %v", err).Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func Ps(w http.ResponseWriter, r *http.Request) {
	spec, err := service.Load(r.PathValue("name"), r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("load: %v", err).Error()))
		return
	}

	res, err := service.Registry[spec.Driver].Ps(context.Background(), spec)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("%s: %v", spec.Driver, err).Error()))
		return
	}

	out, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("json: %v", err).Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func Get(w http.ResponseWriter, r *http.Request) {
	spec, err := service.Load(r.PathValue("name"), r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("load: %v", err).Error()))
		return
	}

	conf, err := service.Get(spec)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("get: %v", err).Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(conf)
}

func Set(w http.ResponseWriter, r *http.Request) {
	conf, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	spec, err := service.Load(r.PathValue("name"), "")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("load: %v", err).Error()))
		return
	}

	if err := service.Set(spec, conf); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("set: %v", err).Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func Run(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var req RunRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("body: %v", err).Error()))
		return
	}

	spec, err := service.Load(r.PathValue("name"), "")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("load: %v", err).Error()))
		return
	}

	fork, err := service.Run(spec, []byte(req.Diff), req.Message)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("run: %v", err).Error()))
		return
	}

	res := RunResponse{IID: fork.IID.Short()}
	raw, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(raw)
}
