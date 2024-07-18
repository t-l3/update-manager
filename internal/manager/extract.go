package manager

import (
	"os"

	"github.com/h2non/filetype"
)

func Extract() {

}

func DetectFiletype(path string) {
	bytes, _ := os.ReadFile(path)
	kind, _ := filetype.Match(bytes)
}
