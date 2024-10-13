package cli

import (
	"context"
	"sync"
	"time"

	"github.com/flohansen/chronos/internal/metric"
)

type Scraper interface {
	Scrape(ctx context.Context, url string) ([]metric.Metric, error)
}

type Storage interface {
	Write(m metric.Metric) error
}

type App struct {
	config  *Config
	scraper Scraper
	storage Storage
}

func NewApp(scraper Scraper, storage Storage, config *Config) *App {
	return &App{
		config:  config,
		scraper: scraper,
		storage: storage,
	}
}

func (a *App) Run(ctx context.Context) <-chan error {
	var wg sync.WaitGroup
	errs := make(chan error)

	for _, target := range a.config.Targets {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case <-time.After(target.Timeout):
					metrics, err := a.scraper.Scrape(ctx, target.URL)
					if err != nil {
						errs <- err
						continue
					}

					for _, m := range metrics {
						if err := a.storage.Write(m); err != nil {
							errs <- err
							continue
						}
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	return errs
}
