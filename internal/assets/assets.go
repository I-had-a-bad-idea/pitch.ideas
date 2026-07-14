package assets

import (
	"embed"
	"io/fs"
	"net/http"
)


//go:embed public/*
var publicFs embed.FS

func Handler() http.Handler {
	sub, err := fs.Sub(publicFs, "public")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(sub))
}