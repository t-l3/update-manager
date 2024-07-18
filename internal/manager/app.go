package manager

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/t-l3/update-manager/internal/config"

	"github.com/0xAX/notificator"
)

func UpdateApplication(app *config.App) {
	logger := log.New(os.Stdout, fmt.Sprintf("app-manager-%s  ", app.Name), log.Ldate|log.Ltime|log.Lmsgprefix)
	logger.Printf("Checking '%s'...", app.Name)

	notify := notificator.New(notificator.Options{
		DefaultIcon: app.Icon,
		AppName:     "update-manager",
	})

	installedCheckProc := exec.Command("/usr/bin/bash", "-c", app.VersioningChecks.Installed)
	installedCheckOutput, err := installedCheckProc.Output()
	if err != nil {
		logger.Println("Failed to run installed version check process\n", err)
		return
	}
	installedVersionString := strings.Trim(string(installedCheckOutput), " \t\n")

	latestCheckProc := exec.Command("/usr/bin/bash", "-c", app.VersioningChecks.Latest)
	latestCheckOutput, err := latestCheckProc.Output()
	if err != nil {
		logger.Println("Failed to run latest version check process\n", err)
		return
	}
	latestVersionString := strings.Trim(string(latestCheckOutput), " \t\n")

	if installedVersionString == latestVersionString {
		logger.Printf("The latest version of %s is currently installed", app.Name)
		return
	} else {
		logger.Printf("Update for %s found\n", app.Name)
		notify.Push("Update found", fmt.Sprintf("Installing update for %s", app.Name), app.Icon, notificator.UR_NORMAL)
	}
}
