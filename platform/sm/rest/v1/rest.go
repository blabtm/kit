package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/blabtm/v2k/platform/sm/service"

	"sigs.k8s.io/yaml"
)

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
	spec, err := service.Load(r.PathValue("name"))

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

func Up(w http.ResponseWriter, r *http.Request) {
	spec, err := service.Load(r.PathValue("name"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("load: %v", err).Error()))
		return
	}

	if err := service.Registry[spec.Driver].Up(context.Background(), spec); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("%s: %v", spec.Driver, err).Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Down(w http.ResponseWriter, r *http.Request) {
	spec, err := service.Load(r.PathValue("name"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("load: %v", err).Error()))
		return
	}

	if err := service.Registry[spec.Driver].Down(context.Background(), spec); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("%s: %v", spec.Driver, err).Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetConfig(w http.ResponseWriter, r *http.Request) {
	yconf, err := service.GetConfig(r.PathValue("name"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("get: %v", err).Error()))
		return
	}

	jconf, err := yaml.YAMLToJSON(yconf)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("encoding: %v", err).Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jconf)
}

func PutConfig(w http.ResponseWriter, r *http.Request) {
	jconf, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	yconf, err := yaml.JSONToYAML(jconf)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Errorf("encoding: %v", err).Error()))
		return
	}

	if err := service.PutConfig(r.PathValue("name"), yconf); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Errorf("put: %v", err).Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
