package api

import (
	"encoding/json"
	"net/http"

	"github.com/daniel-oliveiravas/heartbeat-service/business/heartbeat"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type HandlerConfig struct {
	Logger *zap.Logger
}

type HeartbeatHandler struct {
	heartbeatUsecase *heartbeat.Usecase
}

func New(cfg HandlerConfig) (http.Handler, error) {
	//TODO: Add repository and publisher to heartbeat
	heartbeatUsecase := heartbeat.NewUsecase(cfg.Logger, nil, nil)

	handler := &HeartbeatHandler{
		heartbeatUsecase: heartbeatUsecase,
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

	err := h.heartbeatUsecase.Beat(ctx, beat.toUsecase(identifier))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to read body"))
	}

	w.WriteHeader(200)
}
