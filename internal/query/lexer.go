package query

type TokenType int

const (
	Unknown TokenType = iota
	Literal
	NumberLiteral
)

type Token struct {
	Type  TokenType
	Value string
}

type SimpleLexer struct {
	i     int
	input string
	token Token
}

func NewSimpleLexer(input string) *SimpleLexer {
	return &SimpleLexer{
		i:     0,
		input: input,
	}
}

func (s *SimpleLexer) Next() bool {
	if s.i >= len(s.input) {
		return false
	}

	c := rune(s.input[s.i])

	switch true {
	case isLetter(c):
		s.token = s.parseLiteral()
		return true
	case isNumeric(c):
		s.token = s.parseNumberLiteral()
		return true
	default:
		s.token = Token{Type: Unknown, Value: string(c)}
		s.i++
		return true
	}
}

func (s *SimpleLexer) parseLiteral() Token {
	value := string(s.input[s.i])
	s.i++

	for {
		if s.i >= len(s.input) {
			break
		}

		c := rune(s.input[s.i])

		if isLetter(c) || isNumeric(c) || c == '_' {
			value += string(c)
			s.i++
		} else {
			break
		}
	}

	return Token{Type: Literal, Value: value}
}

func (s *SimpleLexer) parseNumberLiteral() Token {
	value := string(s.input[s.i])
	s.i++

	for {
		if s.i >= len(s.input) {
			break
		}

		c := rune(s.input[s.i])

		if isNumeric(c) {
			value += string(c)
			s.i++
		} else {
			break
		}
	}

	return Token{Type: NumberLiteral, Value: value}
}

func (s *SimpleLexer) Token() Token {
	return s.token
}

func isLetter(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func isNumeric(c rune) bool {
	return c >= '0' && c <= '9'
}
