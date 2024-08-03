package manager

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/h2non/filetype"
	"github.com/t-l3/update-manager/internal/notifications"
)

func (m *Manager) Extract(in string, out string) error {
	kind := DetectFiletype(in)

	var err error

	switch kind { // TODO add zip, 7z, rpm and deb extraction
	case "application/gzip":
		err = m.ExtractGzip(in, out)
	case "application/x-tar":
		err = m.ExtractTar(in, out)
	default:
		m.logger.Println("Finished extracting")
	}
	return err
}

func DetectFiletype(path string) string {
	bytes, _ := os.ReadFile(path)
	kind, _ := filetype.Match(bytes)
	return kind.MIME.Value
}

func (m *Manager) ExtractGzip(in string, out string) error {
	m.logger.Println("Extracting gzip...")
	file, _ := os.Open(in)
	fileStat, _ := file.Stat()
	gzReader, _ := gzip.NewReader(file)

	buffer := make([]byte, 4096)
	tmpFileName := fmt.Sprintf("%s.tmp", in)
	tmpFile, _ := os.Create(tmpFileName)
	total := fileStat.Size()
	read := int64(0)

	notif := notifications.New("msg", m.app.Icon)
	notif.SetInfoMessage(fmt.Sprintf("Extracting %s (Gzip)", m.app.Name))

	for {
		n, err := io.ReadFull(gzReader, buffer)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF { // TODO handle EOF properly
			m.logger.Printf("Failed to extract %s", tmpFileName)
			notif.Terminate(fmt.Sprintf("Extraction of %s failed", m.app.Name))
			return err
		}
		if n == 0 {
			break
		}

		tmpFile.Write(buffer[:n])
		read += int64(n)
		percent := (float64(read) / float64(total)) * float64(100)
		notif.SetPercent(int(percent))
	}
	notif.Terminate(fmt.Sprintf("%s (Gzip) extracted successfully", m.app.Name))
	err := m.Extract(tmpFileName, out)
	os.RemoveAll(tmpFileName)
	return err
}

func (m *Manager) ExtractTar(in string, out string) error {
	m.logger.Println("Extracting tar...")
	file, _ := os.Open(in)
	fileStat, _ := file.Stat()
	tarReader := tar.NewReader(file)

	written := int64(0)
	total := fileStat.Size()

	notif := notifications.New("msg", m.app.Icon)
	notif.SetInfoMessage(fmt.Sprintf("Extracting %s (Tar)", m.app.Name))

	for {
		tarHeader, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error while extracting tar archive %s\n", in)
			notif.Terminate(fmt.Sprintf("Extraction of %s failed", m.app.Name))
			return err
		}

		prefixPath := func(path string) string {
			return fmt.Sprintf("%s/%s", out, path)
		}

		switch tarHeader.Typeflag {
		case tar.TypeDir:
			outPath := prefixPath(tarHeader.Name)
			err := os.MkdirAll(outPath, 0775)
			if err != nil {
				fmt.Printf("Failed to extract directory to '%s'\n", outPath)
				notif.Terminate(fmt.Sprintf("Extraction of %s failed", m.app.Name))
				return err
			}
		case tar.TypeReg:
			outPath := prefixPath(tarHeader.Name)
			newFile, _ := os.Create(outPath)
			n, err := io.Copy(newFile, tarReader)

			written += int64(n)
			percent := (float64(written) / float64(total)) * float64(100)
			notif.SetPercent(int(percent))

			if err != nil {
				fmt.Printf("Failed to extract file to '%s'\n", outPath)
				notif.Terminate(fmt.Sprintf("Extraction of %s failed", m.app.Name))
				return err
			}
			newFile.Close()
		default:
			fmt.Printf("Unable to extract entry in tar archive '%s' of type '%s'\n", tarHeader.Name, string(tarHeader.Typeflag))
			notif.Terminate(fmt.Sprintf("Extraction of %s failed", m.app.Name))
			return errors.New("unhandled tar filetype")
		}
	}

	bytes, _ := io.ReadAll(tarReader)
	tmpFile := fmt.Sprintf("%s.tmp", in)
	os.WriteFile(tmpFile, bytes, 0700)
	notif.Terminate(fmt.Sprintf("%s (Tar) extracted successfully", m.app.Name))
	err := m.Extract(tmpFile, out)
	os.RemoveAll(tmpFile)
	return err
}
