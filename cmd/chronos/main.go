package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	httpv1 "github.com/flohansen/chronos/internal/http"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	scraper := httpv1.NewScraper(&http.Client{})
	log.Fatalf("scraper error: %s", scraper.Start(ctx))
}
