package config

import (
	"os"
)

type RunEnvironment string

const (
	Development RunEnvironment = "dev"
	Production  RunEnvironment = "prod"

	defaultLocalPort = ":8081"
)

var (
	env  RunEnvironment
	port string
)

func GetRunEnvironment() RunEnvironment {

	if env != "" {
		return env
	}

	if os.Getenv("ENV") == "prod" || os.Getenv("ENV") == "production" {
		env = Production
	} else {
		env = Development
	}

	return env
}

// GetPort returns port prepended with `:`
func GetPort() string {
	if port != "" {
		return port
	}

	portNum := os.Getenv("PORT")
	if portNum != "" {
		port = ":" + portNum
		return port
	}

	port = defaultLocalPort
	return port
}
