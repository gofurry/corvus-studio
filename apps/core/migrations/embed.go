package migrations

import (
	"embed"
	"io/fs"
)

//go:embed *.sql
var files embed.FS

func Files() fs.FS {
	return files
}
