package query

type AST struct {
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

func (p *SimpleParser) Parse() (AST, error) {
	for p.lexer.Next() {
		_ = p.lexer.Token()
	}

	return AST{}, nil
}
