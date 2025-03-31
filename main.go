package main

import (
	"embed"
	"flag"
	"html/template"
	"io"
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
	Templates embed.FS

	//go:embed assets/*
	Assets embed.FS
)

func main() {

	port := flag.String("port", "5000", "port to listen on")
	devMode := flag.Bool("dev", false, "run in dev mode")
	flag.Parse()

	var (
		tmplFunc cosmo.ExecuteTemplateFunc
		assets   http.FileSystem
	)

	funcs := template.FuncMap{
		"increment": func(i int) int { return i + 1 },
	}

	if *devMode {
		slog.Info("running in dev mode")
		tmplFunc = func(wr io.Writer, name string, data any) error {
			tmpl, err := template.New("").Funcs(funcs).ParseGlob("templates/*.gohtml")
			if err != nil {
				return err
			}
			return tmpl.ExecuteTemplate(wr, name, data)
		}
		assets = http.Dir(".")
	} else {
		tmpl, err := template.New("").Funcs(funcs).ParseFS(Templates, "templates/*.gohtml")
		if err != nil {
			slog.Error("failed to parse templates", slog.Any("error", err))
			os.Exit(-1)
		}
		tmplFunc = tmpl.ExecuteTemplate
		assets = http.FS(Assets)
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	server := cosmo.NewServer(
		cosmo.FormatBuildVersion(version, commit, buildTime),
		*port,
		httpClient,
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
