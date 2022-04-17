package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Karagar/cyprusExercise/pkg/structs"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	log   *zap.SugaredLogger
	db    *gorm.DB
	route structs.Route
	ctx   context.Context
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	if h.route.TimeoutSec != 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(h.route.TimeoutSec)*time.Second)
	}
	h.ctx = ctx
	defer cancel()

	sqlDB, err := h.db.DB()
	h.handleProblems(w, err)

	err = sqlDB.PingContext(ctx)
	h.handleProblems(w, err)

	if h.route.IsNeedAuth {
		h.log.Info("Stub for authorization")
	}

	if h.route.IsCheckIP {
		h.log.Info("Stub for ip checking")
	}

	if h.route.IsUseQueue {
		h.log.Info("Stub for using a queue")
	}

	if h.route.IsApi {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
	}

	handlerFunc := getHandlerFunc(h.route.Handler)
	handlerFunc(h, w, r)
}

type handlerFunc func(h *Handler, w http.ResponseWriter, r *http.Request)

func getHandlerFunc(funcName string) handlerFunc {
	targetFunc := map[string]handlerFunc{
		"getCompanyHandler":    getCompanyHandler,
		"postCompanyHandler":   postCompanyHandler,
		"putCompanyHandler":    putCompanyHandler,
		"deleteCompanyHandler": deleteCompanyHandler,
	}
	return targetFunc[funcName]
}

func getCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {
	h.log.Info("Serve ", h.route.Path.URL, " (", h.route.Path.Method, ")")
	company := []*structs.Company{}

	result := h.db.WithContext(h.ctx).Table("company").Find(&company)
	h.handleProblems(w, result.Error)

	for _, v := range company {
		v.Guid = handleUuid(v.Uuid)
	}
	body, err := json.Marshal(company)
	h.handleProblems(w, err)
	http.ServeContent(w, r, "index.json", time.Time{}, bytes.NewReader(body))
}

func postCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {}

func putCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {}

func deleteCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {}

func (h *Handler) handleProblems(w http.ResponseWriter, err error) {
	if err != nil {
		h.log.Error(err)
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	}
}

func handleUuid(uuid []byte) string {
	return fmt.Sprintf("%X-%X-%X-%X-%X", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
