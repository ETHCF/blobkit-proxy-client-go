package blobkit_proxy

import "fmt"

var nonRetryableErrors = map[string]struct{}{
	"INVALID_REQUEST":       {},
	"PAYMENT_INVALID":       {},
	"JOB_ALREADY_COMPLETED": {},
	"JOB_EXPIRED":           {},
	"BLOB_TOO_LARGE":        {},
	"SIGNATURE_INVALID":     {},
	// "BLOB_EXECUTION_FAILED": {},
	// "TRANSACTION_FAILED":   {},
}

type Error struct {
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Message)
}

func (e Error) IsRetryable() bool {
	_, ok := nonRetryableErrors[e.Err]
	return !ok
}

// Ensure Error implements the error interface
var _ error = Error{}

type StatusResponse struct {
	Status          string `json:"status"`
	Version         string `json:"version"`
	ChainID         int    `json:"chainId"`
	EscrowContract  string `json:"escrowContract"`
	ProxyFeePercent int    `json:"proxyFeePercent"`
	MaxBlobSize     int    `json:"maxBlobSize"`
	Uptime          int    `json:"uptime"`
	Signer          string `json:"signer"`
}

type BlobWriteRequest struct {
	JobID         string `json:"jobId"`
	PaymentTxHash string `json:"paymentTxHash"`
	Payload       string `json:"payload"`
	Signature     string `json:"signature"`
	Meta          struct {
		AppID       string   `json:"appId"`
		Codec       string   `json:"codec"`
		ContentHash string   `json:"contentHash,omitempty"`
		TTLBlocks   int      `json:"ttlBlocks,omitempty"`
		Timestamp   int64    `json:"timestamp,omitempty"`
		Filename    string   `json:"filename,omitempty"`
		ContentType string   `json:"contentType,omitempty"`
		Tags        []string `json:"tags,omitempty"`
		CallbackURL string   `json:"callbackUrl,omitempty"`
	} `json:"meta"`
}

type BlobWriteResponse struct {
	Success          bool   `json:"success"`
	BlobTxHash       string `json:"blobTxHash"`
	BlockNumber      int    `json:"blockNumber"`
	BlobHash         string `json:"blobHash"`
	Commitment       string `json:"commitment"`
	Proof            string `json:"proof"`
	BlobIndex        int    `json:"blobIndex"`
	CompletionTxHash string `json:"completionTxHash"`
	JobID            string `json:"jobId"`
}
