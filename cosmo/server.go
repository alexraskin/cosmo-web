package cosmo

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/alexraskin/cosmo-web/internal/aws"
)

type ExecuteTemplateFunc func(wr io.Writer, name string, data any) error

type Server struct {
	version     string
	port        string
	server      *http.Server
	assets      http.FileSystem
	tmplFunc    ExecuteTemplateFunc
	awsClient   *aws.Client
	imageConfig ImageConfig
}

func NewServer(version string, port string, assets http.FileSystem, tmplFunc ExecuteTemplateFunc, awsClient *aws.Client, imageConfig ImageConfig) *Server {

	s := &Server{
		version:     version,
		port:        port,
		assets:      assets,
		tmplFunc:    tmplFunc,
		awsClient:   awsClient,
		imageConfig: imageConfig,
	}

	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: s.Routes(),
	}

	return s
}

func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Error while listening", slog.Any("err", err))
		os.Exit(-1)
	}
}

func (s *Server) Close() {
	if err := s.server.Close(); err != nil {
		slog.Error("Error while closing server", slog.Any("err", err))
	}
}

func FormatBuildVersion(version string, commit string, buildTime string) string {
	if len(commit) > 7 {
		commit = commit[:7]
	}

	buildTimeStr := "unknown"
	if buildTime != "unknown" {
		parsedTime, _ := time.Parse(time.RFC3339, buildTime)
		if !parsedTime.IsZero() {
			buildTimeStr = parsedTime.Format(time.ANSIC)
		}
	}
	return fmt.Sprintf("Go Version: %s\nVersion: %s\nCommit: %s\nBuild Time: %s\nOS/Arch: %s/%s\n", runtime.Version(), version, commit, buildTimeStr, runtime.GOOS, runtime.GOARCH)
}
