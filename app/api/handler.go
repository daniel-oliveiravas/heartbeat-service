package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type HandlerConfig struct {
	HeartbeatUsecase HeartbeatUsecase
	Logger           *zap.Logger
}

//go:generate mockery --name=HeartbeatUsecase --filename=heartbeat_usecase.go
type HeartbeatUsecase interface {
	Beat(ctx context.Context, heartbeat heartbeat.Heartbeat) error
	HeartbeatByID(ctx context.Context, id string) (heartbeat.Heartbeat, error)
}

type HeartbeatHandler struct {
	cfg HandlerConfig
}

func New(cfg HandlerConfig) (http.Handler, error) {
	if cfg.HeartbeatUsecase == nil {
		return nil, fmt.Errorf("missing configuration")
	}

	handler := &HeartbeatHandler{
		cfg: cfg,
	}

	router := httprouter.New()
	//TODO: Get id from JWT
	router.POST("/heartbeat/:id", handler.heartbeatSignal)
	router.GET("/heartbeat/:id", handler.getHeartbeatStatus)
	return router, nil
}

func (h *HeartbeatHandler) heartbeatSignal(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("id")
	defer r.Body.Close()

	var beat HeartbeatSignal
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&beat); err != nil {
		h.cfg.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to read body"))
		return
	}

	err := h.cfg.HeartbeatUsecase.Beat(ctx, beat.toUsecase(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to set heartbeat"))
		return
	}

	w.WriteHeader(200)
}

func (h *HeartbeatHandler) getHeartbeatStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	id := ps.ByName("id")
	defer r.Body.Close()

	beat, err := h.cfg.HeartbeatUsecase.HeartbeatByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to get heartbeat"))
		return
	}

	beatJSON, err := json.Marshal(beat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to send response"))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(beatJSON)
}
