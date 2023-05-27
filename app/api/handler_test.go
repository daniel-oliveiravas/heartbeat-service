package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/daniel-oliveiravas/heartbeat-service/app/api"
	"github.com/daniel-oliveiravas/heartbeat-service/app/api/mocks"
	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_HeartbeatEndpoint_ShouldReturnStatusOK(t *testing.T) {
	usecase := mocks.NewHeartbeatUsecase(t)
	handler, err := api.New(api.HandlerConfig{
		HeartbeatUsecase: usecase,
	})
	require.NoError(t, err)

	expectedBeat := heartbeat.Heartbeat{
		ID:        uuid.New().String(),
		Status:    "online",
		Timestamp: time.Now().UTC(),
	}

	usecase.On("Beat", mock.Anything, expectedBeat).Return(nil).Once()

	server := httptest.NewServer(handler)

	httpClient := server.Client()

	beatBody := api.HeartbeatSignal{
		Status:    expectedBeat.Status,
		Timestamp: expectedBeat.Timestamp,
	}

	requestBytes, err := json.Marshal(beatBody)
	require.NoError(t, err)
	body := bytes.NewBuffer(requestBytes)

	resp, err := httpClient.Post(fmt.Sprintf("%s/heartbeat/%s", server.URL, expectedBeat.ID), "application/json", body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
