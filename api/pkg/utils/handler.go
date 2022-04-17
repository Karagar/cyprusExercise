package utils

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Karagar/cyprusExercise/pkg/structs"
	"go.uber.org/zap"
)

type Handler struct {
	log   *zap.SugaredLogger
	db    *sql.DB
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

	err := h.db.PingContext(ctx)
	if err != nil {
		h.log.Error("Cannot ping database")
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	}

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
	// h.log.Info(h.route)
	result := ""
	err := h.db.QueryRowContext(h.ctx, "SELECT [company_name] FROM [exercise].[dbo].[company] where code = 'FirstTestCompanyCode'").Scan(&result)
	if err != nil {
		panic("Cannot select from database")
	}
	fmt.Fprint(w, "{\"data\":\"", result, "\"}")

	// inputQuery := r.URL.Query()
	// inputQuery.Del("contragent_id")
	// r.URL.RawQuery = inputQuery.Encode()

	// for i, back := range b.pipeline {
	// 	if b.ignoreError {
	// 		data.Status = 0
	// 	}
	// 	if data.Status == 200 || data.Status == 0 || i == len(b.pipeline)-1 {
	// 		if back.Serve(w, r, &data) {
	// 			b.log.Info("Backend: Serve HTTP pipeline finished")
	// 			return
	// 		}
	// 	}
	// }
}
func postCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request)   {}
func putCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request)    {}
func deleteCompanyHandler(h *Handler, w http.ResponseWriter, r *http.Request) {}
