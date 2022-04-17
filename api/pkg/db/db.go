package db

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Karagar/cyprusExercise/pkg/utils"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var db *gorm.DB

// New - singleton DB connection wrapper
func New() *gorm.DB {
	if db != nil {
		return db
	}

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
	utils.PanicOnErr(err)

	database := os.Getenv("DB_TITLE")
	if database == "" {
		database = "exercise"
	}

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, password, server, port, database)
	dbConnection, err := gorm.Open(sqlserver.Open(connString), &gorm.Config{})
	utils.PanicOnErr(err)

	db = dbConnection
	return db
}
