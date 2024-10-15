package scrape_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/flohansen/chronos/internal/metric"
	"github.com/flohansen/chronos/internal/metric/scrape"
	"github.com/flohansen/chronos/internal/metric/scrape/mocks"
	"github.com/stretchr/testify/assert"
)

func TestScraper_Scrape(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mocks.NewMockHttpClient(ctrl)

	t.Run("should scrape target", func(t *testing.T) {
		// given
		ctx := context.Background()
		scraper := scrape.NewScraper(client)

		client.EXPECT().
			Do(gomock.Any()).
			Return(&http.Response{
				Body: io.NopCloser(bytes.NewReader([]byte(`
				test_metric1 1234.0
				test_metric2 5678.0
				`))),
			}, nil).
			Times(1)

		// when
		metrics, err := scraper.Scrape(ctx, "any url")

		// then
		assert.Nil(t, err)
		assert.Equal(t, []metric.Metric{
			{Name: "test_metric1", Value: 1234.0},
			{Name: "test_metric2", Value: 5678.0},
		}, metrics)
	})
}
