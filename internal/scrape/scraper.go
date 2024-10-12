package scrape

//go:generate mockgen -destination mocks/scraper.go -package mocks -source scraper.go

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Scraper struct {
	events  chan *http.Response
	client  HttpClient
	targets []HttpTarget
}

func NewScraper(client HttpClient, opts ...ScraperOpt) *Scraper {
	s := &Scraper{
		events: make(chan *http.Response),
		client: client,
	}
	for _, opt := range opts {
		opt(s)
	}

	return s
}

var (
	ErrScrapeCanceled = errors.New("context canceled")
	ErrScrapeHttp     = errors.New("scrape error")
)

// Starts scraping until the context has been canceled.
func (s *Scraper) Start(ctx context.Context) error {
	errs := make(chan error)
	defer close(errs)

	var wg sync.WaitGroup

	for _, target := range s.targets {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				res, err := s.scrapeTarget(ctx, target)
				if err == ErrScrapeCanceled {
					return
				}

				s.events <- res
			}
		}()
	}

	wg.Wait()
	return ErrScrapeCanceled
}

// Scrapes a target using HTTP protocol. If the context has been canceled, it
// returns ErrScrapeCanceled. If HTTP errors occur it returns ErrScrapeHttp.
func (s *Scraper) scrapeTarget(ctx context.Context, target HttpTarget) (*http.Response, error) {
	select {
	case <-ctx.Done():
		return nil, ErrScrapeCanceled
	case <-time.After(target.Timeout):
		url := fmt.Sprintf("%s://%s%s", target.Scheme, target.Host, target.Path)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, ErrScrapeHttp
		}

		res, err := s.client.Do(req)
		if err != nil {
			return nil, ErrScrapeHttp
		}

		return res, nil
	}
}

func (s *Scraper) Events() <-chan *http.Response {
	return s.events
}

type ScraperOpt func(*Scraper)

func WithTarget(target HttpTarget) ScraperOpt {
	return func(s *Scraper) {
		s.targets = append(s.targets, target)
	}
}

func WithConfig(name string) ScraperOpt {
	return func(s *Scraper) {
		f, err := os.Open(name)
		if err != nil {
			fmt.Printf("could not open config: %s\n", err)
			os.Exit(1)
		}
		defer f.Close()

		var cfg Config
		if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
			fmt.Printf("could not yaml decode config: %s\n", err)
			os.Exit(1)
		}

		for _, target := range cfg.Targets {
			WithTarget(target)(s)
		}
	}
}
