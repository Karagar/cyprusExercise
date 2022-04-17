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
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	if h.route.TimeoutSec != 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(h.route.TimeoutSec)*time.Second)
	}
	defer cancel()

	err := h.db.PingContext(ctx)
	if err != nil {
		h.log.Error("Cannot ping database")
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	}

	h.log.Info(h.route)
	result := ""
	// err = conn.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
	err = h.db.QueryRowContext(ctx, "SELECT [company_name] FROM [exercise].[dbo].[company] where code = 'FirstTestCompanyCode'").Scan(&result)
	if err != nil {
		panic("Cannot select from database")
	}
	fmt.Fprint(w, result)

}
