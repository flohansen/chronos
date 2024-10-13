package metric_test

import (
	"bytes"
	"testing"

	"github.com/flohansen/chronos/internal/metric"
	"github.com/stretchr/testify/assert"
)

func TestDecoder_Decode(t *testing.T) {
	t.Run("should return error if a scanned line has invalid syntax", func(t *testing.T) {
		// given
		r := bytes.NewReader([]byte("metric_without_value"))
		d := metric.NewDecoder(r)

		// when
		metrics, err := d.Decode()

		// then
		assert.Error(t, err)
		assert.Equal(t, metric.ErrInvalidSyntax, err)
		assert.Empty(t, metrics)
	})

	t.Run("should return error if a value could not be parsed as float32", func(t *testing.T) {
		// given
		r := bytes.NewReader([]byte("metric_name value_which_is_not_float"))
		d := metric.NewDecoder(r)

		// when
		metrics, err := d.Decode()

		// then
		assert.Error(t, err)
		assert.Equal(t, metric.ErrParseValue, err)
		assert.Empty(t, metrics)
	})

	t.Run("should return metrics while skipping empty lines and ignoring whitespaces", func(t *testing.T) {
		// given
		r := bytes.NewReader([]byte(`
		metric_1 0.0


		metric_2 10.0
				metric_3 20.0
		`))
		d := metric.NewDecoder(r)

		// when
		metrics, err := d.Decode()

		// then
		assert.NoError(t, err)
		assert.Len(t, metrics, 3)
		assert.Equal(t, []metric.Metric{
			{Name: "metric_1", Value: 0.0},
			{Name: "metric_2", Value: 10.0},
			{Name: "metric_3", Value: 20.0},
		}, metrics)
	})

}
