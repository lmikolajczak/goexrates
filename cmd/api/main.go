package main

import (
	"flag"
	"os"

	"github.com/Luqqk/goexrates/internal/data"
	"github.com/Luqqk/goexrates/internal/jsonlog"
)

// The application version number.
const version = "1.0.0"

func main() {
	// Declare an instance of the config struct.
	var cfg config

	// Server conf
	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// Database and connection pool conf
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("EXRATES_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	// Initialize a new jsonlog.Logger which writes any messages *at or above* the INFO
	// severity level to the standard out stream.
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Initialize database connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	// Declare an instance of the application struct,
	// containing all necessary dependencies.
	app := &application{
		config: cfg,
		logger: logger,
		db:     db,
		models: data.NewModels(db),
	}

	app.exposeMetrics()
	// Start the server.
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
