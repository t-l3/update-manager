package manager

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/t-l3/update-manager/internal/config"

	"github.com/0xAX/notificator"
)

func UpdateApplication(app *config.App, tmpPath *string) {
	logger := log.New(os.Stdout, fmt.Sprintf("app-manager-%s  ", app.Name), log.Ldate|log.Ltime|log.Lmsgprefix)
	logger.Printf("Checking '%s'...", app.Name)

	notify := notificator.New(notificator.Options{
		DefaultIcon: app.Icon,
		AppName:     "update-manager",
	})

	_, err := os.ReadDir(app.InstallDir.Path)
	if err == nil {
		installedCheckProc := exec.Command("/usr/bin/bash", "-c", app.VersioningChecks.Installed)
		installedCheckOutput, _ := installedCheckProc.Output()
		installedVersionString := strings.Trim(string(installedCheckOutput), " \t\n")

		latestCheckProc := exec.Command("/usr/bin/bash", "-c", app.VersioningChecks.Latest)
		latestCheckOutput, _ := latestCheckProc.Output()
		latestVersionString := strings.Trim(string(latestCheckOutput), " \t\n")

		if installedVersionString == latestVersionString {
			logger.Printf("The latest version of %s is currently installed", app.Name)
			return
		} else {
			logger.Printf("Update for %s found\n", app.Name)
			notify.Push("Update found", fmt.Sprintf("Installing update for %s", app.Name), app.Icon, notificator.UR_NORMAL)
		}
	}

	logger.Printf("Downloading new version of %s from %s\n", app.Name, app.DownloadUrl)
	res, err := http.Get(app.DownloadUrl)
	if err != nil {
		log.Printf("Failed to download %s", app.Name)
		return
	}

	download, _ := io.ReadAll(res.Body)
	err = os.WriteFile(fmt.Sprintf("%s/%s", *tmpPath, app.Name), download, 0600)
	if err != nil {
		log.Printf("Could not write %s update to temp download directory (%s)", app.Name, app.InstallDir.Path)
	}

	logger.Printf("Removing existing install (%s)\n", app.InstallDir.Path)
	os.RemoveAll(app.InstallDir.Path)
}
