package manager

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/h2non/filetype"
)

func (m *Manager) Extract(in string, out string) {
	kind := DetectFiletype(in)

	switch kind { // TODO add zip, 7z, rpm and deb extraction
	case "application/gzip":
		m.ExtractGzip(in, out)
	case "application/x-tar":
		m.ExtractTar(in, out)
	default:
		m.logger.Println("Finished extracting")
	}
}

func DetectFiletype(path string) string {
	bytes, _ := os.ReadFile(path)
	kind, _ := filetype.Match(bytes)
	return kind.MIME.Value
}

func (m *Manager) ExtractGzip(in string, out string) {
	m.logger.Println("Extracting gzip...")
	file, _ := os.Open(in)
	gzReader, _ := gzip.NewReader(file)
	bytes, _ := io.ReadAll(gzReader)
	tmpFile := fmt.Sprintf("%s.tmp", in)
	os.WriteFile(tmpFile, bytes, 0700)
	m.Extract(tmpFile, out)
	os.RemoveAll(tmpFile)
}

func (m *Manager) ExtractTar(in string, out string) {
	m.logger.Println("Extracting tar...")
	file, _ := os.Open(in)
	tarReader := tar.NewReader(file)

	for {
		tarHeader, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Error while extracting tar archive %s\n", in)
			return
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
				return
			}
		case tar.TypeReg:
			outPath := prefixPath(tarHeader.Name)
			newFile, _ := os.Create(outPath)
			_, err := io.Copy(newFile, tarReader)
			if err != nil {
				fmt.Printf("Failed to extract file to '%s'\n", outPath)
				return
			}
			newFile.Close()
		default:
			fmt.Printf("Unable to extract entry in tar archive '%s' of type '%s'\n", tarHeader.Name, string(tarHeader.Typeflag))
			return
		}
	}

	bytes, _ := io.ReadAll(tarReader)
	tmpFile := fmt.Sprintf("%s.tmp", in)
	os.WriteFile(tmpFile, bytes, 0700)
	m.Extract(tmpFile, out)
	os.RemoveAll(tmpFile)
}
