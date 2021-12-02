package config

import (
	"os"
)

type RunEnvironment string

const (
	Development RunEnvironment = "dev"
	Production  RunEnvironment = "prod"

	defaultLocalPort = ":8081"

	defaultProjectID = "crypto-lost-chances"

	lostChancesCalcLocalhost = "localhost:8082"
)

var (
	env       RunEnvironment
	port      string
	projectID string

	lostChancesCalcHost string
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

func GetLostChancesCalcHost() string {
	if lostChancesCalcHost != "" {
		return lostChancesCalcHost
	}

	lostChancesCalcHost = os.Getenv("LOST_CHANCES_CALC_HOST")
	if lostChancesCalcHost == "" {
		lostChancesCalcHost = lostChancesCalcLocalhost
	}

	return lostChancesCalcHost
}
