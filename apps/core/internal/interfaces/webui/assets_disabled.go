//go:build !corvus_webui

package webui

import "io/fs"

func Files() (fs.FS, error) {
	return nil, nil
}
