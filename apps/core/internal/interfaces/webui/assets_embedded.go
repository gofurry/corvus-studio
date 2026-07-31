//go:build corvus_webui

package webui

import (
	"embed"
	"io/fs"
)

//go:embed dist
var embedded embed.FS

func Files() (fs.FS, error) {
	return fs.Sub(embedded, "dist")
}
