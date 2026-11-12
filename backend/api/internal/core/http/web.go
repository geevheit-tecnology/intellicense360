package http

import (
	"embed"
	"io/fs"
	stdhttp "net/http"
)

//go:embed web/index.html web/src
var embeddedWebFS embed.FS

func loadEmbeddedWeb() (indexHTML []byte, staticHandler stdhttp.Handler) {
	sub, err := fs.Sub(embeddedWebFS, "web")
	if err != nil {
		return nil, nil
	}
	indexHTML, _ = fs.ReadFile(sub, "index.html")
	staticHandler = stdhttp.FileServer(stdhttp.FS(sub))
	return indexHTML, staticHandler
}

