package config

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	// The build environment (used to select a config)
	Environment string `envconfig:"ENV_NAME" default:"dev"`
	// The name of the application
	AppName string `ignored:"true"`
	// The current version number of the app
	AppVersion string `ignored:"true"`
	// Shorthand app name
	AppShortName string `ignored:"true"`
	// The URL to expose the application on
	AppAddress string `envconfig:"ADDR" default:":8080"`
}

var configuration Configuration
var onceConfig = sync.Once{}

// Returns the config object. Only runs once
func GetConfiguration() *Configuration {
	onceConfig.Do(func() {
		log.Print("loading environment variables")
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
		configuration = Configuration{
			AppName:      "Skeleton",
			AppShortName: "SKEL",
			AppVersion:   "1.0.0",
		}
		err = envconfig.Process(configuration.AppShortName, &configuration)
		if err != nil {
			log.Fatal(err)
		}
	})
	return &configuration
}
