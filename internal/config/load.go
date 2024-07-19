package config

import (
	"flag"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig() AppConfig {
	configFile := flag.String("config", "/etc/update-manager/config.yaml", "relative or absolute file path to update-manager's config")

	configFileHandle, err := os.Open(*configFile)
	if err != nil {
		log.Fatal("Error encountered while opening config file", err)
	}

	configFileBytes, err := io.ReadAll(configFileHandle)
	if err != nil {
		log.Fatal("IO Error encountered while reading config file", err)
	}
	configFileHandle.Close()

	appConfig := AppConfig{
		TmpDownloadLocation: "/tmp/update-manager-download",
		SystrayIcon:         "/etc/update-manager/icons/update-manager.png",
	}
	err = yaml.Unmarshal(configFileBytes, &appConfig)
	if err != nil {
		log.Fatal("Cannot parse config file.\n", err)
	}

	if appConfig.TmpDownloadLocation == "/" {
		log.Fatal("Cannot use '/' as tmp directory")
	}

	return appConfig
}
