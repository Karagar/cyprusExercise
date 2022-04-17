package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Karagar/cyprusExercise/pkg/structs"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	Log   *zap.SugaredLogger
	DB    *gorm.DB
	Route structs.Route
	ctx   context.Context
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	if h.Route.TimeoutSec != 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(h.Route.TimeoutSec)*time.Second)
	}
	h.ctx = ctx
	defer cancel()

	sqlDB, err := h.DB.DB()
	if err != nil {
		h.handleProblems(w, err)
		return
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		h.handleProblems(w, err)
		return
	}

	if h.Route.IsNeedAuth {
		h.Log.Info("Stub for authorization")
	}

	if h.Route.IsCheckIP {
		h.Log.Info("Stub for ip checking")
	}

	if h.Route.IsUseQueue {
		h.Log.Info("Stub for using a queue")
	}

	if h.Route.IsApi {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
	}

	handlerFunc := getHandlerFunc(h.Route.Handler)
	handlerFunc(h, w, r)
}

func (h *Handler) handleProblems(w http.ResponseWriter, err error) {
	h.Log.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "{\"error\":", strconv.Quote(err.Error()), "}")
}
