package blobkit_proxy

import "time"

const (
	MainnetDefaultProxyURL = "https://proxy.blobkit.org"
	DefaultRequestTimeout  = 60 * time.Second
	DefaultRetryDelay      = 3 * time.Second
)

type ProxyConfig struct {
	Endpoint   string
	Timeout    time.Duration
	MaxRetries int
	RetryDelay time.Duration
}

func (pc ProxyConfig) WithDefaults() ProxyConfig {
	if pc.Endpoint == "" {
		pc.Endpoint = MainnetDefaultProxyURL
	}
	if pc.Timeout == 0 {
		pc.Timeout = DefaultRequestTimeout
	}
	if pc.MaxRetries == 0 {
		pc.MaxRetries = 3
	}
	if pc.RetryDelay == 0 {
		pc.RetryDelay = DefaultRetryDelay
	}
	return pc
}
