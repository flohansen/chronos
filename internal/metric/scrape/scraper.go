package scrape

//go:generate mockgen -destination mocks/scraper.go -package mocks -source scraper.go

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	var metrics []metric.Metric
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("invalid metric syntax: '%s'", line)
		}

		value, err := strconv.ParseFloat(tokens[1], 32)
		if err != nil {
			return nil, fmt.Errorf("could not parse value: %s", err)
		}

		metrics = append(metrics, metric.Metric{
			Name:  strings.TrimSpace(tokens[0]),
			Value: float32(value),
		})
	}

	return metrics, nil
}
