package config

import "os"

type RunEnvironment string

const (
	Development RunEnvironment = "dev"
	Production  RunEnvironment = "prod"
)

func GetRunEnvironment() RunEnvironment {
	if os.Getenv("ENV") == "prod" || os.Getenv("ENV") == "production" {
		return Production
	}

	return Development
}
