package http

//go:generate mockgen -destination mocks/scraper.go -package mocks -source scraper.go

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpTarget struct {
	Scheme  string
	Host    string
	Path    string
	Timeout time.Duration
}

type Scraper struct {
	client  HttpClient
	targets []HttpTarget
}

func NewScraper(client HttpClient) *Scraper {
	return &Scraper{
		client: client,
	}
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
				err := s.scrapeTarget(ctx, target)
				if err == ErrScrapeCanceled {
					return
				}
			}
		}()
	}

	wg.Wait()
	return ErrScrapeCanceled
}

// Scrapes a target using HTTP protocol. If the context has been canceled, it
// returns ErrScrapeCanceled. If HTTP errors occur it returns ErrScrapeHttp.
func (s *Scraper) scrapeTarget(ctx context.Context, target HttpTarget) error {
	select {
	case <-ctx.Done():
		return ErrScrapeCanceled
	case <-time.After(target.Timeout):
		url := fmt.Sprintf("%s://%s%s", target.Scheme, target.Host, target.Path)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return ErrScrapeHttp
		}

		_, err = s.client.Do(req)
		if err != nil {
			return ErrScrapeHttp
		}

		return nil
	}
}
