package scrape_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/flohansen/chronos/internal/scrape"
	"github.com/flohansen/chronos/internal/scrape/mocks"
	"github.com/stretchr/testify/assert"
)

func TestScraper_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mocks.NewMockHttpClient(ctrl)

	t.Run("should scrape target two times using interval", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())
		scraper := scrape.NewScraper(
			client,
			scrape.WithTarget(scrape.HttpTarget{
				Scheme:  "http",
				Host:    "example.com",
				Path:    "/metrics",
				Timeout: 10 * time.Millisecond,
			}),
		)

		client.EXPECT().
			Do(gomock.Any()).
			Return(&http.Response{}, nil).
			Times(1)

		client.EXPECT().
			Do(gomock.Any()).
			Do(func(r *http.Request) {
				cancel()
			}).
			Return(&http.Response{}, nil).
			Times(1)

		// when
		t1 := time.Now()
		err := scraper.Start(ctx)

		// then
		assert.NotNil(t, err)
		assert.Equal(t, "context canceled", err.Error())
		assert.GreaterOrEqual(t, time.Since(t1), 20*time.Millisecond)
	})

	t.Run("should not stop when http errors occur", func(t *testing.T) {
		// given
		ctx, cancel := context.WithCancel(context.Background())
		scraper := scrape.NewScraper(
			client,
			scrape.WithTarget(scrape.HttpTarget{
				Scheme:  "http",
				Host:    "example.com",
				Path:    "/metrics",
				Timeout: 10 * time.Millisecond,
			}),
		)

		client.EXPECT().
			Do(gomock.Any()).
			Return(nil, errors.New("some error")).
			Times(2)

		client.EXPECT().
			Do(gomock.Any()).
			Do(func(r *http.Request) {
				cancel()
			}).
			Return(&http.Response{}, nil).
			Times(1)

		// when
		t1 := time.Now()
		err := scraper.Start(ctx)

		// then
		assert.NotNil(t, err)
		assert.Equal(t, "context canceled", err.Error())
		assert.GreaterOrEqual(t, time.Since(t1), 20*time.Millisecond)
	})
}
