package manager

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/0xAX/notificator"
)

func (m *Manager) InstallApp() {
	downloadPath := fmt.Sprintf("%s/%s", *m.tmpDir, m.app.Name)
	extractPath := fmt.Sprintf("%s-dir", downloadPath)

	m.Extract(downloadPath, extractPath)

	m.logger.Printf("Removing existing install (%s)\n", m.app.InstallDir.Path)
	os.RemoveAll(m.app.InstallDir.Path)

	m.logger.Printf("Moving extracted files from %s to %s\n", extractPath, m.app.InstallDir.Path)
	err := os.Rename(extractPath, m.app.InstallDir.Path)
	if err != nil {
		cmd := exec.Command("cp", "-r", extractPath, m.app.InstallDir.Path)
		cmd.Start()
		err := cmd.Wait()
		if err != nil {
			m.logger.Printf("Failed to copy from %s to %s", extractPath, m.app.InstallDir.Path)
		}
		os.RemoveAll(extractPath)
		os.RemoveAll(downloadPath)
	}

	if len(m.app.PostInstallScript) > 0 {
		exec.Command("/usr/bin/bash", "-c", m.app.PostInstallScript).Start()
	}

	m.notify.Push("Update complete", fmt.Sprintf("Completed update of %s to %s", m.app.Name, m.app.InstallDir.Path), m.app.Icon, notificator.UR_NORMAL)

}
