package manager

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/0xAX/notificator"
)

func (m *Manager) CheckVersion() bool {
	m.logger.Printf("Checking '%s'...", m.app.Name)

	_, err := os.ReadDir(m.app.InstallDir.Path)
	if err == nil {
		installedCheckProc := exec.Command("/usr/bin/bash", "-c", m.app.VersioningChecks.Installed)
		installedCheckOutput, _ := installedCheckProc.Output()
		installedVersionString := strings.Trim(string(installedCheckOutput), " \t\n")

		latestCheckProc := exec.Command("/usr/bin/bash", "-c", m.app.VersioningChecks.Latest)
		latestCheckOutput, _ := latestCheckProc.Output()
		latestVersionString := strings.Trim(string(latestCheckOutput), " \t\n")

		if installedVersionString == latestVersionString {
			m.logger.Printf("The latest version of %s is currently installed", m.app.Name)
			return false
		} else {
			m.logger.Printf("Update for %s found\n", m.app.Name)
			m.notify.Push("Update found", fmt.Sprintf("Installing update for %s", m.app.Name), m.app.Icon, notificator.UR_NORMAL)
			return true
		}
	} else {
		m.logger.Printf("Install path for %s is not present. Continuing download\n", m.app.Name)
		m.notify.Push("Installing app", fmt.Sprintf("Installing %s", m.app.Name), m.app.Icon, notificator.UR_NORMAL)
		return true
	}
}
