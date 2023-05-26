package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/daniel-oliveiravas/heartbeat-service/app/api"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HeartbeatEndpoint_ShouldReturnStatusOK(t *testing.T) {
	handler, err := api.New(api.HandlerConfig{})
	require.NoError(t, err)

	server := httptest.NewServer(handler)

	httpClient := server.Client()

	id := uuid.New().String()
	resp, err := httpClient.Post(fmt.Sprintf("%s/heartbeat/%s", server.URL, id), "application/json", nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
