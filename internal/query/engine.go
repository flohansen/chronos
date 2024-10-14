package query

import (
	"fmt"

	"github.com/flohansen/chronos/internal/metric"
)

type Parser interface {
	Parse() (AST, error)
}

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Exec(query string) ([]metric.Metric, error) {
	parser := NewSimpleParser(NewSimpleLexer(query))

	_, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("could not parse query: %s", err)
	}

	return nil, nil
}
