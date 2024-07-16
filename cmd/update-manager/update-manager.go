package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/t-l3/update-manager/internal/config"

	"gopkg.in/yaml.v2"
)

func main() {
	configFile := flag.String("config", "/etc/update-manager/config.yaml", "relative or absolute file path to update-manager's config")

	configFileHandle, err := os.Open(*configFile)
	if err != nil {
		log.Fatal("Error encountered while opening config file", err)
	}

	configFileBytes, err := io.ReadAll(configFileHandle)
	if err != nil {
		log.Fatal("IO Error encountered while reading config file", err)
	}

	appConfig := config.AppConfig{}
	err = yaml.Unmarshal(configFileBytes, &appConfig)
	if err != nil {
		log.Fatal("Cannot parse config yaml", err)
	}
}
