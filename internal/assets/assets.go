package assets

import (
	"embed"
	"io/fs"
	"net/http"
)


//go:embed static/*
var StaticFs embed.FS

func Handler() http.Handler {
	sub, err := fs.Sub(StaticFs, "static")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(sub))
}