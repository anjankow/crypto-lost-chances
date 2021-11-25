package config

import (
	"os"
)

type RunEnvironment string

const (
	Development RunEnvironment = "dev"
	Production  RunEnvironment = "prod"

	defaultLocalPort  = ":8081"
	defaultProdDomain = "crypto-lost-chances.appspot.com"
	localDomain       = "localhost" + defaultLocalPort
)

var (
	env    RunEnvironment
	port   string
	domain string
)

func GetRunEnvironment() RunEnvironment {

	if env == "" {
		if os.Getenv("ENV") == "prod" || os.Getenv("ENV") == "production" {
			env = Production
		} else {
			env = Development
		}
	}

	return env
}

// GetPort returns port prepended with `:`
func GetPort() string {
	if port == "" {
		if GetRunEnvironment() == Production {
			port = ":" + os.Getenv("PORT")
		} else {
			port = defaultLocalPort
		}
	}

	return port
}

func GetDomainAddr() string {
	if domain != "" {
		return domain
	}

	if GetRunEnvironment() != Production {
		domain = localDomain
		return domain
	}

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		domain = defaultProdDomain
	} else {
		domain = apiUrl
	}

	return domain
}
