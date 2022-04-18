package handler

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

// ServeHTTP is for HTTP request handling
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
		h.handleProblems(w, http.StatusInternalServerError, err)
		return
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		h.handleProblems(w, http.StatusBadGateway, err)
		return
	}

	if h.Route.IsNeedAuth {
		token := url.QueryEscape(utils.ReadCookie("token", r))
		if token == "" {
			h.handleProblems(w, http.StatusForbidden, errors.New("Unauthorized"))
			return
		}
		claims, err := utils.ReadJwt(token)
		if err != nil {
			h.handleProblems(w, http.StatusForbidden, err)
			return
		}
		if claims["name"] != "John Doe" {
			h.handleProblems(w, http.StatusForbidden, errors.New("Unauthorized"))
			return
		}
	}

	if h.Route.IsCheckIP {
		ipapiClient := http.Client{}
		clientIP := utils.ReadUserIP(r, false)
		if clientIP == "" {
			h.handleProblems(w, http.StatusForbidden, errors.New("Sorry, You've got wrong placement"))
			return
		}
		checkUrl := fmt.Sprintf("https://ipapi.co/%s/country/", clientIP)
		req, err := http.NewRequest("GET", checkUrl, nil)
		if err != nil {
			h.handleProblems(w, http.StatusBadGateway, err)
			return
		}
		req.Header.Set("User-Agent", "ipapi.co/#go-v1.5")
		resp, err := ipapiClient.Do(req)
		if err != nil {
			h.handleProblems(w, http.StatusBadGateway, err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			h.handleProblems(w, http.StatusBadGateway, err)
			return
		}
		if string(body) != "CY" {
			h.handleProblems(w, http.StatusForbidden, errors.New("Sorry, You've got wrong placement"))
			return
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

// handleProblems is an error wrapper
// log error and return error status and description to client
func (h *Handler) handleProblems(w http.ResponseWriter, status int, err error) {
	h.Log.Error(err)
	w.WriteHeader(status)
	fmt.Fprint(w, "{\"error\":", strconv.Quote(err.Error()), "}")
}
