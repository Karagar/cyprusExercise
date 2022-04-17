package handler

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Karagar/cyprusExercise/pkg/structs"
	"github.com/Karagar/cyprusExercise/pkg/utils"
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
		ipapiClient := http.Client{}
		clientIP := utils.ReadUserIP(r, false)
		checkUrl := fmt.Sprintf("https://ipapi.co/%s/country/", clientIP)
		req, err := http.NewRequest("GET", checkUrl, nil)
		req.Header.Set("User-Agent", "ipapi.co/#go-v1.5")
		resp, err := ipapiClient.Do(req)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.handleProblems(w, err)
		}
		if string(body) != "CY" {
			h.handleProblems(w, errors.New("Sorry, You've got wrong placement"))
		}
	}

	if h.Route.IsUseQueue {
		h.Log.Info("Stub for using a queue")
	}

	if h.Route.IsApi {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
	}

	handlerFunc := getHandlerFunc(h.Route.Handler)
	h.Log.Info("Serve ", h.Route.Path.URL, " (", h.Route.Path.Method, ")")
	handlerFunc(h, w, r)
}

func (h *Handler) handleProblems(w http.ResponseWriter, err error) {
	h.Log.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "{\"error\":", strconv.Quote(err.Error()), "}")
}
