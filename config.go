package blobkit_proxy

import "time"

type ProxyConfig struct {
	Endpoint string
	Timeout  time.Duration
}

type StatusResponse struct {
	Status          string `json:"status"`
	Version         string `json:"version"`
	ChainId         int    `json:"chainId"`
	EscrowContract  string `json:"escrowContract"`
	ProxyFeePercent int    `json:"proxyFeePercent"`
	MaxBlobSize     int    `json:"maxBlobSize"`
	Uptime          int    `json:"uptime"`
	Signer          string `json:"signer"`
}
