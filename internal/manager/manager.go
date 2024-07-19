package manager

import (
	"fmt"
	"log"

	"github.com/0xAX/notificator"
	"github.com/t-l3/update-manager/internal/config"
)

type Manager struct {
	app          *config.App
	tmpDir       *string
	logger       *log.Logger
	notify       *notificator.Notificator
	downloadPath *string
	extractPath  *string
}

func New(app *config.App, tmpDir *string, logger *log.Logger, notify *notificator.Notificator) Manager {
	downloadPath := fmt.Sprintf("%s/%s", *tmpDir, app.Name)
	extractPath := fmt.Sprintf("%s-dir", downloadPath)

	if downloadPath == "/" {
		log.Fatalf("Cannot use '/' as download path (%s)", app.Name)
	}

	return Manager{
		app:          app,
		tmpDir:       tmpDir,
		logger:       logger,
		notify:       notify,
		downloadPath: &downloadPath,
		extractPath:  &extractPath,
	}
}
