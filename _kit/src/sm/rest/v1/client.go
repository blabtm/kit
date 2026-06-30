package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/blabtm/v2k/platform/sm/service"
)

type Client struct {
	Base string
}

func NewClient(host string, port int) *Client {
	return &Client{
		Base: fmt.Sprintf("http://%s:%d/api", host, port),
	}
}

func (c *Client) Ls() ([]string, error) {
	res, err := http.Get(c.Base)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", string(raw))
	}

	var services []string

	if err := json.Unmarshal(raw, &services); err != nil {
		return nil, fmt.Errorf("body: json: %w", err)
	}

	return services, nil
}

func (c *Client) Ps(name string, iid string) (*service.Status, error) {
	params := url.Values{}

	if iid != "" {
		params.Add("id", iid)
	}

	res, err := http.Get(c.Base + "/" + name + "/ps" + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", string(raw))
	}

	var status service.Status

	if err := json.Unmarshal(raw, &status); err != nil {
		return nil, fmt.Errorf("body: json: %w", err)
	}

	return &status, nil
}

func (c *Client) Get(name string, iid string) ([]byte, error) {
	params := url.Values{}

	if iid != "" {
		params.Add("id", iid)
	}

	res, err := http.Get(c.Base + "/" + name + "?" + params.Encode())
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

	req, err := http.NewRequest(http.MethodPut, c.Base+"/"+name, body)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer res.Body.Close()

	raw, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("body: %w", err)
	}

	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("%s", string(raw))
	}

	return nil
}

func (c *Client) Fork(name string, raw []byte, msg string) (string, error) {
	runReq := ForkRequest{
		Message: msg,
		Diff:    string(raw),
	}

	raw, err := json.Marshal(runReq)
	if err != nil {
		return "", fmt.Errorf("body: %w", err)
	}

	body := bytes.NewReader(raw)

	req, err := http.NewRequest(http.MethodPut, c.Base+"/"+name+"/fork", body)
	if err != nil {
		return "", fmt.Errorf("http: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http: %w", err)
	}
	defer res.Body.Close()

	raw, err = io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("body: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("%s", string(raw))
	}

	var response ForkResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return "", fmt.Errorf("response: %w", err)
	}

	return response.IID, nil
}
