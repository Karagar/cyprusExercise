package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Karagar/cyprusExercise/pkg/structs"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// serverStruct - structure for the core server entities
type ServerStruct struct {
	config       *structs.Config
	logger       *zap.SugaredLogger
	dbConnection *gorm.DB
	address      string
}

func (s *ServerStruct) Declare() {
	s.getConfig()
	s.getLogger()
	s.getDBConnection()
	s.getAddress()
}

// getConfig - singleton config wrapper
func (s *ServerStruct) getConfig() {
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "../scripts/config.json"
	}

	if s.config != nil {
		return
	}

	s.config = &structs.Config{}
	configContent := MustReadFile(configPath)
	panicOnErr(json.Unmarshal(configContent, s.config))
}

// getLogger - singleton SugarLogger wrapper
func (s *ServerStruct) getLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "DEBUG"
	}

	if s.logger != nil {
		return
	}

	if logLevel == "DEBUG" {
		s.logger = getDevelopmentLog()
	} else {
		s.logger = getProductionLog()
	}
}

// getDBConnection - singleton DB connection wrapper
func (s *ServerStruct) getDBConnection() {
	server := os.Getenv("DB_SERVER")
	if server == "" {
		server = "cyprus_db"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "sa"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "Newer_use_it_in_prod111"
	}

	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		portStr = "1433"
	}
	port, err := strconv.Atoi(portStr)
	panicOnErr(err)

	database := os.Getenv("DB_TITLE")
	if database == "" {
		database = "exercise"
	}

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, password, server, port, database)
	s.dbConnection, err = gorm.Open(sqlserver.Open(connString), &gorm.Config{})
	panicOnErr(err)
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
	panicOnErr(err)
	defer sqlDB.Close()
	r := mux.NewRouter()

	for _, route := range s.config.Routes {
		handler := &Handler{
			log:   s.logger,
			db:    s.dbConnection,
			route: route,
		}
		r.Path(route.Path.URL).Methods(route.Path.Method).Handler(handler)
	}

	s.logger.Info("Start: Listen - ", s.address)
	s.logger.Fatal(http.ListenAndServe(s.address, r))
}
