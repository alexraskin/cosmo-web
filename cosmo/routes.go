package cosmo

import (
	"io"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Use(httprate.Limit(
		100,
		time.Minute,
		httprate.WithLimitHandler(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
			}),
		),
	))

	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/assets/")

		file, err := s.assets.Open(path)
		if err != nil {
			slog.Error("asset not found", slog.String("path", path), slog.Any("error", err))
			http.Error(w, "Asset not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			slog.Error("error getting file stat", slog.String("path", path), slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}

		http.ServeContent(w, r, path, stat.ModTime(), file.(io.ReadSeeker))
	})

	r.Handle("/robots.txt", serveFile(s.assets, "robots.txt"))
	r.Handle("/sitemap.xml", serveFile(s.assets, "sitemap.xml"))

	r.Get("/", s.index)

	r.NotFound(s.notFound)

	return r
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	data := DefaultGalleryConfig()

	err := s.tmplFunc(w, "index.gohtml", data)
	if err != nil {
		slog.Error("template execution failed", slog.Any("error", err))
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

func serveFile(fs http.FileSystem, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := fs.Open(path)
		if err != nil {
			slog.Error("file not found", slog.String("path", path), slog.Any("error", err))
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()
		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		_, _ = io.Copy(w, file)
	}
}
