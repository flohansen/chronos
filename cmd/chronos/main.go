package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"

	"github.com/flohansen/chronos/internal/cli"
	"github.com/flohansen/chronos/internal/metric/scrape"
	"github.com/flohansen/chronos/internal/metric/storage"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	config := cli.NewConfig(cli.FromFile("chronos.yml"))

	if err := os.MkdirAll(config.Storage.Directory, 0777); err != nil {
		log.Fatalf("could create storage directory: %s", err)
	}

	filename := path.Join(config.Storage.Directory, "current")
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("could not open file: %s", err)
	}
	defer f.Close()

	scraper := scrape.NewScraper(&http.Client{})
	storage := storage.NewFileWriter(f)
	app := cli.NewApp(scraper, storage, config)

	for err := range app.Run(ctx) {
		log.Printf("error: %s", err)
	}
}
