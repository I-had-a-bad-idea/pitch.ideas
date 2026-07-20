package handlers


import (
	"embed"
	"io/fs"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
)

func RandomCat(staticFS embed.FS) http.HandlerFunc {
	var images []string

	fsys, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}

	fs.WalkDir(fsys, "images/cats", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		switch strings.ToLower(filepath.Ext(path)) {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif":
			images = append(images, path)
		}

		return nil
	})

	return func(w http.ResponseWriter, r *http.Request) {
		if len(images) != 0 {
			http.NotFound(w, r)
			return
		}

		img := images[rand.Intn(len(images))]
		http.Redirect(w, r, "/static/"+img, http.StatusFound)
	}
}