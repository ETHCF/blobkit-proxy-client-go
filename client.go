package blobkit_proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/cockroachdb/errors"
)

type Client struct {
	conf   ProxyConfig
	client *http.Client
}

func NewClient(conf ProxyConfig) *Client {
	return &Client{
		conf:   conf,
		client: &http.Client{Timeout: conf.Timeout},
	}
}

func (c *Client) get(ctx context.Context, method, path string, out any) error {
	ep, err := url.JoinPath(c.conf.Endpoint, path)
	if err != nil {
		return errors.Wrap(err, "failed to create endpoint URL")
	}

	req, err := http.NewRequestWithContext(ctx, method, ep, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d error: %s", resp.StatusCode, string(data))
	}

	if err := json.Unmarshal(data, out); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetStatus(ctx context.Context) (StatusResponse, error) {
	var status StatusResponse
	if err := c.get(ctx, http.MethodGet, "/api/v1/health", &status); err != nil {
		return StatusResponse{}, err
	}
	return status, nil
}
