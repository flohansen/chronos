package query_test

import (
	"testing"

	"github.com/flohansen/chronos/internal/query"
	"github.com/stretchr/testify/assert"
)

func TestParseMetric(t *testing.T) {
	// given
	input := "metric_name@12345=10.0"

	// when
	m, err := query.ParseMetric(input)

	// then
	assert.NoError(t, err)
	assert.Equal(t, query.MetricRow{
		Name:      "metric_name",
		Value:     10.0,
		Timestamp: 12345,
	}, m)
}
