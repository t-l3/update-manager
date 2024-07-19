package manager

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func (m *Manager) DownloadApp() {
	downloadPath := fmt.Sprintf("%s/%s", *m.tmpDir, m.app.Name)

	_, err := os.Open(downloadPath)

	if err != nil {
		m.logger.Printf("Downloading new version of %s from %s\n", m.app.Name, m.app.DownloadUrl)
		res, err := http.Get(m.app.DownloadUrl)

		if err != nil {
			m.logger.Printf("Failed to download %s", m.app.Name)
			return
		}

		download, _ := io.ReadAll(res.Body)
		err = os.WriteFile(downloadPath, download, 0600)

		if err != nil {
			m.logger.Printf("Could not write %s update to temp download directory (%s)", m.app.Name, m.app.InstallDir.Path)
		}
	} else {
		m.logger.Printf("Download already present, reusing existing download for %s", m.app.Name) // TODO Add file size check to confirm presence of data
	}
}
