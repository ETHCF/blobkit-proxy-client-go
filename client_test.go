package blobkit_proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetStatus(t *testing.T) {
	t.Parallel()
	client := NewClient(ProxyConfig{})

	status, err := client.GetStatus(t.Context())
	assert.NoError(t, err, "GetStatus should not return an error")

	assert.Equal(t, "healthy", status.Status, "unexpected status")
}

func Test_GetEscrowContract(t *testing.T) {
	t.Parallel()
	client := NewClient(ProxyConfig{})

	escrowContract, err := client.GetEscrowContract(t.Context())
	assert.NoError(t, err, "GetEscrowContract should not return an error")
	assert.Len(t, escrowContract, 42, "escrow contract should be a valid address length")
}
