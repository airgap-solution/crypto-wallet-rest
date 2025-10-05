package service_test

import (
	"net/http"
	"testing"

	"github.com/airgap-solution/crypto-wallet-rest/internal/core/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleError(t *testing.T) {
	t.Parallel()
	resp, err := service.HandleError(assert.AnError)

	require.NoError(t, err)
	require.Equal(t, http.StatusNotImplemented, resp.Code)

	body, ok := resp.Body.(service.Error)
	require.True(t, ok)
	require.Equal(t, assert.AnError.Error(), body.Message)
}
