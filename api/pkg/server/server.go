package server

import (
	"net/http"
	"os"

	"github.com/Karagar/cyprusExercise/pkg/config"
	"github.com/Karagar/cyprusExercise/pkg/db"
	"github.com/Karagar/cyprusExercise/pkg/handler"
	"github.com/Karagar/cyprusExercise/pkg/logger"
	"github.com/Karagar/cyprusExercise/pkg/structs"
	"github.com/Karagar/cyprusExercise/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// serverStruct - structure for the core server entities
type ServerStruct struct {
	config       *structs.Config
	logger       *zap.SugaredLogger
	dbConnection *gorm.DB
	address      string
}

var server *ServerStruct

func New() *ServerStruct {
	if server != nil {
		return server
	}

	server = &ServerStruct{}
	server.getAddress()
	server.config = config.New()
	server.logger = logger.New()
	server.dbConnection = db.New()
	return server
}

// getAddress - func to get listening address
func (s *ServerStruct) getAddress() {
	address := os.Getenv("PORT")
	if address == "" {
		address = ":8080"
	}
	s.address = address
}

// Serve - core handler
func (s *ServerStruct) Serve() {
	sqlDB, err := s.dbConnection.DB()
	utils.PanicOnErr(err)
	defer sqlDB.Close()
	r := mux.NewRouter()

	for _, route := range s.config.Routes {
		handler := &handler.Handler{
			Log:   s.logger,
			DB:    s.dbConnection,
			Route: route,
		}
		r.Path(route.Path.URL).Methods(route.Path.Method).Handler(handler)
	}

	s.logger.Info("Start: Listen - ", s.address)
	s.logger.Fatal(http.ListenAndServe(s.address, r))
}
