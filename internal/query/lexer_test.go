package query_test

import (
	"testing"

	"github.com/flohansen/chronos/internal/query"
	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	// given
	lexer := query.NewSimpleLexer(`metric_name`)

	// when
	var tokens []query.Token
	for lexer.Next() {
		tokens = append(tokens, lexer.Token())
	}

	// then
	assert.Equal(t, []query.Token{
		{Type: query.Literal, Value: "metric_name"},
	}, tokens)
}
