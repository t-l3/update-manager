package manager

import (
	"fmt"
	"log"
	"os"

	"github.com/t-l3/update-manager/internal/config"

	"github.com/0xAX/notificator"
)

func UpdateApplication(app* config.App) {
	logger := log.New(os.Stdout, fmt.Sprintf("app-manager-%s  ", app.Name), log.Ldate | log.Ltime | log.Lmsgprefix)
	logger.Printf("Checking '%s'...", app.Name)

	notify := notificator.New(notificator.Options{
    DefaultIcon: app.Icon,
    AppName:     "update-manager",
  })

  notify.Push("Update found", fmt.Sprintf("Installing update for %s", app.Name), app.Icon, notificator.UR_NORMAL)
}
