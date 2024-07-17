package manager

import (
	"log"

	"github.com/t-l3/update-manager/internal/config"
)

func UpdateApplication(app config.App) {
	log.Printf(" == Checking '%s'... == ", app.Name)
}
