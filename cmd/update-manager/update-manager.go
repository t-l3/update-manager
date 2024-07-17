package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"

	"github.com/t-l3/update-manager/internal/config"
	"github.com/t-l3/update-manager/internal/manager"

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

	appConfig := config.AppConfig{
		TmpDownloadLocation: "/tmp/update-manager-download",
	}
	err = yaml.Unmarshal(configFileBytes, &appConfig)
	if err != nil {
		log.Fatal("Cannot parse config file.\n", err)
	}

	err = os.MkdirAll(appConfig.TmpDownloadLocation, 0775)
	if err != nil {
		log.Fatal("Error while creating download directory", err)
	}

	log.Printf("  === Starting app checks ===  ")

	var wg sync.WaitGroup

	for _, app := range appConfig.Apps {
		wg.Add(1)
		go updateApplication(&app, &wg)
	}

	wg.Wait()

	os.RemoveAll(appConfig.TmpDownloadLocation)
}

func updateApplication(app* config.App, wg* sync.WaitGroup) {
	manager.UpdateApplication(app)
	wg.Done()
}