package v1

import (
	"bytes"
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
		Base: fmt.Sprintf("http://%s:%d/v1", host, port),
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

func (c *Client) Get(name string) ([]byte, error) {
	res, err := http.Get(c.Base + "/" + name + "/get")

	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("body: %w", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(raw))
	}

	return raw, nil
}

func (c *Client) Set(name string, raw []byte) error {
	body := bytes.NewReader(raw)
	req, err := http.NewRequest(http.MethodPut, c.Base+"/"+name+"/set", body)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return fmt.Errorf("http: %w", err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer res.Body.Close()

	raw, err = io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("body: %w", err)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%s", string(raw))
	}

	return nil
}
