package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/blabtm/v2k/platform/sm/service"
)

type Client struct {
	Base string
}

func NewClient(host string, port int) *Client {
	return &Client{
		Base: fmt.Sprintf("http://%s:%d/v1/svc", host, port),
	}
}

func (c *Client) Ls() ([]string, error) {
	res, err := http.Get(c.Base)

	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	} else {
		defer res.Body.Close()
	}

	dat, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("body: %w", err)
	}

	var services []string

	if err := json.Unmarshal(dat, &services); err != nil {
		return nil, fmt.Errorf("body: json: %w", err)
	}

	return services, nil
}

func (c *Client) Ps(name string) (*service.Status, error) {
	res, err := http.Get(c.Base + "/" + name + "/ps")

	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	} else {
		defer res.Body.Close()
	}

	dat, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("body: %w", err)
	}

	var status service.Status

	if err := json.Unmarshal(dat, &status); err != nil {
		return nil, fmt.Errorf("body: json: %w", err)
	}

	return &status, nil
}

func (c *Client) Up(name string) error {
	res, err := http.Get(c.Base + "/" + name + "/up")

	if err != nil {
		return fmt.Errorf("http: %w", err)
	} else {
		defer res.Body.Close()
	}

	return nil
}

func (c *Client) Down(name string) error {
	res, err := http.Get(c.Base + "/" + name + "/down")

	if err != nil {
		return fmt.Errorf("http: %w", err)
	} else {
		defer res.Body.Close()
	}

	return nil
}

