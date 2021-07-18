package main

// Define a config struct to hold all the configuration settings for the application.
// Read these configuration settings from cmd-line flags when the application starts.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}
