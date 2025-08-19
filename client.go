// Package blobkit_proxy provides a client for interacting with the BlobKit proxy service.
package blobkit_proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cockroachdb/errors"
)

// Client represents a BlobKit proxy client that handles HTTP communication
// with the BlobKit proxy service.
type Client struct {
	conf   ProxyConfig
	client *http.Client
}

// NewClient creates a new BlobKit proxy client with the provided configuration.
// The configuration is normalized with default values applied automatically.
func NewClient(conf ProxyConfig) *Client {
	return &Client{
		conf:   conf.WithDefaults(),
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

// GetStatus retrieves the current status and health information from the BlobKit proxy service.
// It returns a StatusResponse containing service health details and configuration.
func (c *Client) GetStatus(ctx context.Context) (StatusResponse, error) {
	var status StatusResponse
	if err := c.get(ctx, http.MethodGet, "/api/v1/health", &status); err != nil {
		return StatusResponse{}, err
	}
	return status, nil
}

// GetEscrowContract retrieves the escrow contract address from the BlobKit proxy service.
// It calls GetStatus internally and extracts the EscrowContract field.
func (c *Client) GetEscrowContract(ctx context.Context) (string, error) {
	status, err := c.GetStatus(ctx)
	if err != nil {
		return "", err
	}
	return status.EscrowContract, nil
}

func (c *Client) writeBlob(ctx context.Context, blobReq BlobWriteRequest) (BlobWriteResponse, error) {
	ep, err := url.JoinPath(c.conf.Endpoint, "/api/v1/blob/write")
	if err != nil {
		return BlobWriteResponse{}, errors.Wrap(err, "failed to create write blob endpoint URL")
	}
	data, err := json.Marshal(blobReq)
	if err != nil {
		return BlobWriteResponse{}, errors.Wrap(err, "failed to marshal request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewBuffer(data))
	if err != nil {
		return BlobWriteResponse{}, errors.Wrap(err, "failed to create request")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return BlobWriteResponse{}, errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return BlobWriteResponse{}, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		var errResp Error
		if err := json.Unmarshal(respData, &errResp); err != nil {
			return BlobWriteResponse{}, fmt.Errorf("unexpected status code: %d error: %s", resp.StatusCode, string(respData))
		}
		return BlobWriteResponse{}, errResp
	}

	var blobResp BlobWriteResponse
	if err := json.Unmarshal(respData, &blobResp); err != nil {
		return BlobWriteResponse{}, errors.Wrap(err, "failed to unmarshal response")
	}

	return blobResp, nil
}

// WriteBlob writes blob data to the BlobKit proxy service with automatic retry logic.
// It retries failed requests up to MaxRetries times with RetryDelay between attempts.
// Non-retryable errors (as determined by the Error.IsRetryable method) will not be retried.
func (c *Client) WriteBlob(ctx context.Context, blobReq BlobWriteRequest) (BlobWriteResponse, error) {

	for i := range c.conf.MaxRetries {

		writeResp, err := c.writeBlob(ctx, blobReq)
		if err == nil {
			return writeResp, nil
		}

		if respError, ok := err.(Error); ok && !respError.IsRetryable() {
			return BlobWriteResponse{}, err
		}

		if i == c.conf.MaxRetries-1 {
			return BlobWriteResponse{}, errors.Wrap(err, "failed to write blob")
		}
		time.Sleep(c.conf.RetryDelay)
	}

	return BlobWriteResponse{}, errors.New("unexpected error in WriteBlob method, should not reach here")
}
