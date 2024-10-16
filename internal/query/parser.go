package query

type Node interface {
	Eval(metric MetricRow) bool
}

type MetricSelector struct {
	Name string
}

func (n *MetricSelector) Eval(data MetricRow) bool {
	return data.Name == n.Name
}

type Lexer interface {
	Next() bool
	Token() Token
}

type SimpleParser struct {
	lexer Lexer
}

func NewSimpleParser(lexer Lexer) *SimpleParser {
	return &SimpleParser{
		lexer: lexer,
	}
}

func (p *SimpleParser) Parse() (Node, error) {
	for p.lexer.Next() {
		token := p.lexer.Token()
		switch token.Type {
		case Literal:
			return &MetricSelector{
				Name: token.Value,
			}, nil
		}
	}

	return nil, nil
}
