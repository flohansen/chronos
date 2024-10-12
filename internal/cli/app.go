package cli

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"

	"github.com/flohansen/chronos/internal/scrape"
	"github.com/flohansen/chronos/internal/storage"
	"gopkg.in/yaml.v3"
)

func Init() error {
	f, err := os.Create("chronos.yml")
	if err != nil {
		return fmt.Errorf("could not open file: %s", err)
	}
	defer f.Close()

	if err := yaml.NewEncoder(f).Encode(scrape.DefaultConfig); err != nil {
		return fmt.Errorf("could not encode default config: %s", err)
	}

	return nil
}

func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	scraper := scrape.NewScraper(&http.Client{}, scrape.WithConfig("chronos.yml"))
	storage := storage.NewFileWriter()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-scraper.Events():
				if err := storage.Write("metric_name", rand.Float32()); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	return fmt.Errorf("scraper error: %s", scraper.Start(ctx))
}
