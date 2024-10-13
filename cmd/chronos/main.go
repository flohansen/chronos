package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/flohansen/chronos/internal/cli"
	"github.com/flohansen/chronos/internal/metric/scrape"
	"github.com/flohansen/chronos/internal/metric/storage"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	config := cli.NewConfig(cli.FromFile("chronos.yml"))
	scraper := scrape.NewScraper(&http.Client{})
	storage := storage.NewFileWriter()
	app := cli.NewApp(scraper, storage, config)

	for err := range app.Run(ctx) {
		log.Printf("error: %s", err)
	}
}
