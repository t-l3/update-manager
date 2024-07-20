package manager

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/t-l3/update-manager/internal/notifications"
)

func (m *Manager) DownloadApp() {
	downloadPath := fmt.Sprintf("%s/%s", *m.tmpDir, m.app.Name)

	_, err := os.Open(downloadPath)

	if err != nil {
		m.logger.Printf("Downloading new version of %s from %s\n", m.app.Name, m.app.DownloadUrl)
		notif := notifications.New("msg", m.app.Icon)
		notif.SetInfoMessage(fmt.Sprintf("Downloading %s", m.app.Name))

		res, err := http.Get(m.app.DownloadUrl)
		if err != nil {
			m.logger.Printf("Failed to download %s", m.app.Name)
			notif.Terminate(fmt.Sprintf("Download of %s failed", m.app.Name))
			return
		}

		total := res.ContentLength
		read := int64(0)

		download, _ := os.Create(downloadPath)
		buffer := make([]byte, 4096)

		for {
			n, err := io.ReadFull(res.Body, buffer)

			if err != nil && err != io.EOF {
				m.logger.Printf("Failed to download %s", m.app.Name)
				notif.Terminate(fmt.Sprintf("Download of %s failed", m.app.Name))
				return
			}
			if n == 0 {
				break
			}

			download.Write(buffer[:n])
			read += int64(n)

			percent := (float64(read) / float64(total)) * float64(100)

			notif.SetPercent(int(percent))
		}

		notif.Terminate(fmt.Sprintf("%s downloaded successfully", m.app.Name))
	} else {
		m.logger.Printf("Download already present, reusing existing download for %s", m.app.Name) // TODO Add file size check to confirm presence of data
	}
}
