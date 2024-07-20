package manager

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/0xAX/notificator"
)

func (m *Manager) InstallApp() {
	m.Extract(*m.downloadPath, *m.extractPath)

	m.logger.Printf("Removing existing install (%s)\n", m.app.InstallDir.Path)
	os.RemoveAll(m.app.InstallDir.Path)

	if len(m.app.PreInstallScript) > 0 {
		exec.Command("/usr/bin/bash", "-c", m.app.PreInstallScript).Start()
	}

	m.logger.Printf("Moving extracted files from %s to %s\n", *m.extractPath, m.app.InstallDir.Path)
	err := os.Rename(*m.extractPath, m.app.InstallDir.Path)
	if err != nil {
		cmd := exec.Command("cp", "-r", *m.extractPath, m.app.InstallDir.Path)
		cmd.Start()
		err := cmd.Wait()
		if err != nil {
			m.logger.Printf("Failed to copy from %s to %s", *m.extractPath, m.app.InstallDir.Path)
		}
		os.RemoveAll(*m.extractPath)
		os.RemoveAll(*m.downloadPath)
	}

	if len(m.app.PostInstallScript) > 0 {
		exec.Command("/usr/bin/bash", "-c", m.app.PostInstallScript).Start()
	}

	m.notify.Push("Update complete", fmt.Sprintf("Completed update of %s to %s", m.app.Name, m.app.InstallDir.Path), m.app.Icon, notificator.UR_NORMAL)

}
