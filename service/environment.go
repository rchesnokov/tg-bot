package service

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// Environment ... list of required env variables' values
type Environment struct {
	DatabaseURL string
	Mode        string
	Token       string
}

type mode string

const (
	modeDevelopment = "dev"
	modeProduction  = "prod"
)

// InitEnviroment ... return Environment instance
func InitEnviroment() *Environment {
	loadEnvFromFile()
	databaseURL := lookupEnvVariable("DB_URL")
	mode := defineAppMode(lookupEnvVariable("ENV"))
	token := lookupEnvVariable("TOKEN")

	return &Environment{
		databaseURL,
		mode,
		token,
	}
}

func loadEnvFromFile() {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func lookupEnvVariable(v string) string {
	value, defined := os.LookupEnv(v)
	if !defined {
		log.Fatalf("%s env variable is not provided", v)
	}

	return value
}

func defineAppMode(envMode string) string {
	m := mode(envMode)

	if m.isValid() {
		return envMode
	}

	return modeProduction
}

func (m *mode) isValid() bool {
	switch *m {
	case modeDevelopment, modeProduction:
		return true
	}

	return false
}
