package service_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/airgap-solution/crypto-wallet-rest/internal/core/service"
	"github.com/stretchr/testify/require"
)

func TestHandleError(t *testing.T) {
	testErr := errors.New("test error message")

	resp, err := service.HandleError(testErr)

	require.NoError(t, err)
	require.Equal(t, http.StatusNotImplemented, resp.Code)

	body, ok := resp.Body.(service.Error)
	require.True(t, ok)
	require.Equal(t, "test error message", body.Message)
}
