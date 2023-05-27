package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
	"github.com/julienschmidt/httprouter"
)

type HandlerConfig struct {
	HeartbeatUsecase HeartbeatUsecase
}

//go:generate mockery --name=HeartbeatUsecase --filename=heartbeat_usecase.go
type HeartbeatUsecase interface {
	Beat(ctx context.Context, heartbeat heartbeat.Heartbeat) error
}

type HeartbeatHandler struct {
	cfg HandlerConfig
}

func New(cfg HandlerConfig) (http.Handler, error) {
	//TODO: Add repository and publisher to heartbeat

	if cfg.HeartbeatUsecase == nil {
		return nil, fmt.Errorf("missing configuration")
	}

	handler := &HeartbeatHandler{
		cfg: cfg,
	}

	router := httprouter.New()
	//TODO: Get id from JWT
	router.POST("/heartbeat/:id", handler.heartbeatSignal)
	return router, nil
}

func (h *HeartbeatHandler) heartbeatSignal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	identifier := ps.ByName("id")
	defer r.Body.Close()

	var beat HeartbeatSignal
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&beat); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to read body"))
	}

	err := h.cfg.HeartbeatUsecase.Beat(ctx, beat.toUsecase(identifier))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to read body"))
	}

	w.WriteHeader(200)
}
