package cli

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/flohansen/chronos/internal/metric"
	"github.com/flohansen/chronos/internal/query"
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- a.startQueryServer(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range a.startScrapers(ctx) {
			errs <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errs)
	}()

	return errs
}

func (a *App) startQueryServer(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /query", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		engine := query.NewEngine()
		rows, err := engine.Exec(string(b))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(rows)
	})

	srv := &http.Server{
		Addr:    ":2020",
		Handler: mux,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	return srv.ListenAndServe()
}

func (a *App) startScrapers(ctx context.Context) <-chan error {
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
				case <-time.After(target.Interval):
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
