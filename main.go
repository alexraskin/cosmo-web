package main

import (
	"embed"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexraskin/cosmo-web/cosmo"
)

var (
	version   = "unknown"
	commit    = "unknown"
	buildTime = "unknown"
)

var (
	//go:embed templates/*
	templates embed.FS

	//go:embed assets/*
	assets embed.FS
)

func main() {

	port := flag.String("port", "5000", "port to listen on")
	flag.Parse()

	funcs := template.FuncMap{
		"increment": func(i int) int { return i + 1 },
	}

	tmpl, err := template.New("").Funcs(funcs).ParseFS(templates, "templates/*.gohtml")
	if err != nil {
		slog.Error("failed to parse templates", slog.Any("error", err))
		os.Exit(-1)
	}
	tmplFunc := tmpl.ExecuteTemplate

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	server := cosmo.NewServer(
		cosmo.FormatBuildVersion(version, commit, buildTime),
		*port,
		httpClient,
		templates,
		assets,
		tmplFunc,
	)

	go server.Start()
	defer server.Close()

	slog.Info("started cosmo-web", slog.Any("listen_addr", *port))
	si := make(chan os.Signal, 1)
	signal.Notify(si, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-si
	slog.Info("shutting down cosmo-web")

}
