package scrape

//go:generate mockgen -destination mocks/scraper.go -package mocks -source scraper.go

import (
	"context"
	"net/http"

	"github.com/flohansen/chronos/internal/metric"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Scraper struct {
	client HttpClient
}

func NewScraper(client HttpClient) *Scraper {
	s := &Scraper{
		client: client,
	}

	return s
}

// Scrapes a HTTP target.
func (s *Scraper) Scrape(ctx context.Context, url string) ([]metric.Metric, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	metrics, err := metric.NewDecoder(res.Body).Decode()
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
