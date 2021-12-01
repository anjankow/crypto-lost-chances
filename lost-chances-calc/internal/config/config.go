package config

import (
	"os"
)

type RunEnvironment string

const (
	Development RunEnvironment = "dev"
	Production  RunEnvironment = "prod"

	defaultLocalPort = ":8082"

	defaultProjectID = "crypto-lost-chances"
)

var (
	env       RunEnvironment
	port      string
	projectID string
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

func GetProjectID() string {
	if projectID != "" {
		return projectID
	}

	projectID = os.Getenv("PROJECT_ID")
	if projectID == "" {
		projectID = defaultProjectID
	}

	return projectID
}
