package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

var conn *sql.DB
var dbPassword = os.Getenv("MSSQL_SA_PASSWORD")

func main() {
	var (
		server   string = "cyprus_db"
		user     string = "sa"
		password string = dbPassword
		database string = "exercise"
		port     int    = 1433
		err      error
	)

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", server, user, password, port, database)
	conn, err = sql.Open("sqlserver", connString)

	if err != nil {
		panic("Cannot connect to database")
	} else {
		fmt.Println("Connected!")
	}
	defer conn.Close()

	http.HandleFunc("/", hp)
	http.ListenAndServe(":8080", nil)
}

func hp(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	err := conn.PingContext(ctx)
	if err != nil {
		panic("Cannot ping database")
	}

	result := ""
	// err = conn.QueryRowContext(ctx, "SELECT @@version").Scan(&result)
	err = conn.QueryRowContext(ctx, "SELECT [company_name] FROM [exercise].[dbo].[company] where code = 'FirstTestCompanyCode'").Scan(&result)
	if err != nil {
		panic("Cannot select from database")
	}
	fmt.Fprint(w, result)
}
