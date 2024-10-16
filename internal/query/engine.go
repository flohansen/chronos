package query

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Parser interface {
	Parse() (Node, error)
}

type Engine struct {
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Exec(query string) ([]MetricRow, error) {
	parser := NewSimpleParser(NewSimpleLexer(query))

	node, err := parser.Parse()
	if err != nil {
		return nil, fmt.Errorf("could not parse query: %s", err)
	}

	f, err := os.Open("data/current")
	if err != nil {
		return nil, fmt.Errorf("could not open file: %s", err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var metrics []MetricRow
	for scanner.Scan() {
		line := scanner.Text()
		m, err := ParseMetric(line)
		if err != nil {
			return nil, fmt.Errorf("could not parse metric: %s", err)
		}

		if node.Eval(m) {
			metrics = append(metrics, m)
		}
	}

	return metrics, nil
}

type MetricRow struct {
	Name      string  `json:"name"`
	Value     float32 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

func ParseMetric(input string) (MetricRow, error) {
	m := MetricRow{}
	lexer := NewSimpleLexer(input)

	if ok := lexer.Next(); !ok {
		return MetricRow{}, errors.New("invalid metric syntax")
	}
	token := lexer.Token()
	if token.Type != Literal {
		return MetricRow{}, errors.New("expected metric name")
	}
	m.Name = token.Value

	if ok := lexer.Next(); !ok {
		return MetricRow{}, errors.New("invalid metric syntax")
	}
	token = lexer.Token()
	if token.Value != "@" {
		return MetricRow{}, errors.New("expected @")
	}

	if ok := lexer.Next(); !ok {
		return MetricRow{}, errors.New("invalid metric syntax")
	}
	token = lexer.Token()
	if token.Type != NumberLiteral {
		return MetricRow{}, errors.New("expected timestamp")
	}

	var err error
	m.Timestamp, err = strconv.ParseInt(token.Value, 10, 64)
	if err != nil {
		return MetricRow{}, fmt.Errorf("could not parse timestamp: %s", err)
	}

	if ok := lexer.Next(); !ok {
		return MetricRow{}, errors.New("invalid metric syntax")
	}
	token = lexer.Token()
	if token.Value != "=" {
		return MetricRow{}, errors.New("expected =")
	}

	if ok := lexer.Next(); !ok {
		return MetricRow{}, errors.New("invalid metric syntax")
	}
	token = lexer.Token()
	if token.Type != NumberLiteral {
		return MetricRow{}, errors.New("expected value")
	}
	value := token.Value

	if lexer.Next() {
		token := lexer.Token()
		if token.Value != "." {
			return MetricRow{}, errors.New("expected .")
		}

		if ok := lexer.Next(); !ok {
			return MetricRow{}, errors.New("invalid metric syntax")
		}
		token = lexer.Token()
		if token.Type != NumberLiteral {
			return MetricRow{}, errors.New("expected floating point number")
		}
		value += "." + token.Value
	}

	valueFloat, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return MetricRow{}, fmt.Errorf("could not parse value: %s", err)
	}

	m.Value = float32(valueFloat)
	return m, nil
}
