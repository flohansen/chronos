package metric

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

var (
	ErrInvalidSyntax = errors.New("invalid syntax")
	ErrParseValue    = errors.New("invalid value")
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

func (d *Decoder) Decode() ([]Metric, error) {
	var metrics []Metric

	scanner := bufio.NewScanner(d.r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		tokens := strings.Split(line, " ")
		if len(tokens) != 2 {
			return nil, ErrInvalidSyntax
		}

		value, err := strconv.ParseFloat(tokens[1], 32)
		if err != nil {
			return nil, ErrParseValue
		}

		metrics = append(metrics, Metric{
			Name:  strings.TrimSpace(tokens[0]),
			Value: float32(value),
		})
	}

	return metrics, nil
}
