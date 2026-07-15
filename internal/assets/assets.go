package assets

import (
	"embed"
	"io/fs"
	"net/http"
)


//go:embed static/*
var staticFs embed.FS

func Handler() http.Handler {
	sub, err := fs.Sub(staticFs, "static")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(sub))
}