package manager

import (
	"log"

	"github.com/0xAX/notificator"
	"github.com/t-l3/update-manager/internal/config"
)

type Manager struct {
	app    *config.App
	tmpDir *string
	logger *log.Logger
	notify *notificator.Notificator
}

func New(app *config.App, tmpDir *string, logger *log.Logger, notify *notificator.Notificator) Manager {
	return Manager{
		app:    app,
		tmpDir: tmpDir,
		logger: logger,
		notify: notify,
	}
}
