package main

import (
	"context"
	"embed"
	"flag"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexraskin/cosmo-web/cosmo"
	"github.com/alexraskin/cosmo-web/internal/aws"
	"github.com/alexraskin/cosmo-web/internal/config"
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
	configFilePath := flag.String("config", "config.yaml", "path to config file")
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

	config, err := config.LoadConfig(*configFilePath)
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		os.Exit(-1)
	}

	awsClient, err := aws.New(context.Background(), aws.Config{
		AccountID: config.AWSConfig.AccountID,
		Bucket:    config.AWSConfig.Bucket,
		Region:    config.AWSConfig.Region,
		AccessKey: config.AWSConfig.AccessKey,
		SecretKey: config.AWSConfig.SecretKey,
	})

	if err != nil {
		slog.Error("failed to create aws client", slog.Any("error", err))
		os.Exit(-1)
	}

	imageConfig := cosmo.ImageConfig{
		BaseURL:         config.ImageConfig.BaseURL,
		ThumbnailParams: config.ImageConfig.ThumbnailParams,
		FullsizeParams:  config.ImageConfig.FullsizeParams,
		Folder:          config.ImageConfig.Folder,
	}

	server := cosmo.NewServer(
		cosmo.FormatBuildVersion(version, commit, buildTime),
		*port,
		assets,
		tmplFunc,
		awsClient,
		imageConfig,
	)

	go server.Start()
	defer server.Close()

	slog.Info("started cosmo-web", slog.Any("listen_addr", *port))
	si := make(chan os.Signal, 1)
	signal.Notify(si, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-si
	slog.Info("shutting down cosmo-web")

}
