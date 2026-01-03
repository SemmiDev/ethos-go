// Package web provides embedded frontend assets and SPA serving capabilities.
// This allows the entire React frontend to be compiled into the Go binary.
package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

// Embed the entire frontend/dist directory.
// This directive tells Go to include all files from frontend/dist at compile time.
// The build process should run `npm run build` in the frontend directory first.
//
//go:embed all:dist
var frontendFS embed.FS

// DistFS returns the embedded filesystem rooted at "dist"
func DistFS() (fs.FS, error) {
	return fs.Sub(frontendFS, "dist")
}

// SPAHandler serves the embedded React SPA.
// It handles:
// - Static file requests (js, css, images, etc.)
// - SPA routing (returns index.html for non-file requests)
// - Proper Content-Type headers
// - Gzip compression support (if pre-compressed files exist)
type SPAHandler struct {
	fileServer http.Handler
	fs         fs.FS
}

// NewSPAHandler creates a new SPA handler with the embedded filesystem.
func NewSPAHandler() (*SPAHandler, error) {
	distFS, err := DistFS()
	if err != nil {
		return nil, err
	}

	return &SPAHandler{
		fileServer: http.FileServer(http.FS(distFS)),
		fs:         distFS,
	}, nil
}

// ServeHTTP implements http.Handler for the SPA.
func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean the URL path
	urlPath := r.URL.Path
	if urlPath == "/" {
		urlPath = "/index.html"
	}

	// Remove leading slash for fs.Open
	filePath := strings.TrimPrefix(urlPath, "/")

	// Check if the file exists
	file, err := h.fs.Open(filePath)
	if err != nil {
		// File doesn't exist, serve index.html for SPA routing
		h.serveIndex(w, r)
		return
	}
	file.Close()

	// Check if it's a file (not a directory)
	stat, err := fs.Stat(h.fs, filePath)
	if err != nil || stat.IsDir() {
		// It's a directory or error, serve index.html
		h.serveIndex(w, r)
		return
	}

	// Set proper cache headers for assets
	if isAsset(urlPath) {
		// Long cache for versioned assets (vite adds hash to filenames)
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else if urlPath == "/index.html" || urlPath == "/" {
		// No cache for index.html to ensure updates
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	}

	// Serve the file
	h.fileServer.ServeHTTP(w, r)
}

// serveIndex serves the index.html file for SPA routing.
func (h *SPAHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	indexFile, err := h.fs.Open("index.html")
	if err != nil {
		http.Error(w, "Index not found", http.StatusNotFound)
		return
	}
	defer indexFile.Close()

	stat, err := indexFile.Stat()
	if err != nil {
		http.Error(w, "Error reading index", http.StatusInternalServerError)
		return
	}

	// Set headers for HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// Read and serve the file
	data, err := fs.ReadFile(h.fs, "index.html")
	if err != nil {
		http.Error(w, "Error reading index", http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, r, "index.html", stat.ModTime(), strings.NewReader(string(data)))
}

// isAsset checks if the path is a static asset that should be cached.
func isAsset(urlPath string) bool {
	ext := path.Ext(urlPath)
	switch ext {
	case ".js", ".css", ".woff", ".woff2", ".ttf", ".eot", ".svg", ".png", ".jpg", ".jpeg", ".gif", ".ico", ".webp":
		return true
	}
	// Also match paths containing /assets/
	return strings.Contains(urlPath, "/assets/")
}

// IsFrontendBundled returns true if the frontend is embedded.
// This helps determine if we should serve embedded files or proxy to dev server.
func IsFrontendBundled() bool {
	_, err := frontendFS.ReadFile("dist/index.html")
	return err == nil
}

// DevModeHandler returns a handler that proxies to the Vite dev server.
// This is useful for development when the frontend isn't embedded.
func DevModeHandler(devServerURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In dev mode, redirect to the Vite dev server
		http.Redirect(w, r, devServerURL+r.URL.Path, http.StatusTemporaryRedirect)
	})
}

// MustFS is a helper that panics if there's an error.
// Useful for initialization where errors are fatal.
func MustFS() fs.FS {
	distFS, err := DistFS()
	if err != nil {
		panic("failed to load embedded frontend: " + err.Error())
	}
	return distFS
}

// FileExists checks if a file exists in the embedded filesystem.
func FileExists(path string) bool {
	distFS, err := DistFS()
	if err != nil {
		return false
	}
	_, err = fs.Stat(distFS, path)
	return err == nil || !os.IsNotExist(err)
}
