package main

import (
	"dss-api/internal/data"
	"dss-api/internal/driver"
	"fmt"
	"log"
	"net/http"
	"os"
)

// config is the type for all aplication configuration
type config struct {
	port int
}

// application is the type for all data
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	// db       *driver.DB
	models      data.Models
	environment string
}

// main is the main entry point of our application
func main() {
	var cfg config
	cfg.port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// dsn := "host=localhost port=5432 user=postgres password=password dbname=dssapi sslmode=disable timezone=UTC connect_timeout=5"
	dsn := os.Getenv("DSN")
	environment := os.Getenv("ENV")
	db, err := driver.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	defer db.SQL.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		// db:       db,
		models:      data.New(db.SQL),
		environment: environment,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {

	app.infoLog.Println("API listening on port", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}
