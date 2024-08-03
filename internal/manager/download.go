package manager

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/t-l3/update-manager/internal/notifications"
)

func (m *Manager) DownloadApp() error {
	downloadPath := fmt.Sprintf("%s/%s", *m.tmpDir, m.app.Name)

	res, err := http.Get(m.app.DownloadUrl)
	if err != nil {
		m.logger.Printf("Failed to send request to download URL for %s", m.app.Name)
		return err
	}
	total := res.ContentLength

	existingFile, err := os.Open(downloadPath)
	if err == nil {
		if stat, err := existingFile.Stat(); err == nil && stat.Size() == total {
			m.logger.Printf("Download of same size already present, reusing existing download for %s", m.app.Name)
			return nil
		}
		m.logger.Printf("Existing download for %s has incorrect file size, re-downloading", m.app.Name)
		os.Remove(downloadPath)
	}

	m.logger.Printf("Downloading new version of %s from %s\n", m.app.Name, m.app.DownloadUrl)
	notif := notifications.New("msg", m.app.Icon)
	notif.SetInfoMessage(fmt.Sprintf("Downloading %s", m.app.Name))

	read := int64(0)

	download, _ := os.Create(downloadPath)
	buffer := make([]byte, 4096)

	for {
		n, err := io.ReadFull(res.Body, buffer)

		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF { // TODO handle EOF properly
			m.logger.Printf("Failed to download %s", m.app.Name)
			notif.Terminate(fmt.Sprintf("Download of %s failed", m.app.Name))
			return err
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

	return nil
}
