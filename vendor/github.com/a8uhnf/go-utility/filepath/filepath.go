package filepath

import (
	"path/filepath"
	"strings"
)

// Filename returns file-name from path.
func Filename(path string) string {
	dir := filepath.Dir(path) + "/"
	return strings.TrimPrefix(path, dir)
}
