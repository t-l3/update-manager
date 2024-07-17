package manager

import (
	"fmt"
	"log"
	"os"

	"github.com/t-l3/update-manager/internal/config"
)

func UpdateApplication(app* config.App) {
	logger := log.New(os.Stdout, fmt.Sprintf("app-manager-%s  ", app.Name), log.Ldate | log.Ltime | log.Lmsgprefix)
	logger.Printf("Checking '%s'...", app.Name)
}
