package parser

import (
	"bytes"
	"fmt"
)

func tokenize(s string) ([]token, error) {
	var buf bytes.Buffer
	tokens := []token{}
	i := 0
	for {
		if i >= len(s) {
			return tokens, nil
		}
		switch s[i] {
		case ' ', '\t', '\r', '\n':
			i++
		case '(', ')', '=', '[', ']', '|':
			tokens = append(tokens, token{
				t:      s[i],
				offset: i,
				length: 1,
			})
			i++
		case '$':
			tok := token{offset: i}
			buf.Reset()
			buf.WriteByte('$')
			i++
			for i < len(s) && symchar(s[i]) {
				buf.WriteByte(s[i])
				i++
			}
			s := buf.String()
			buf.Reset()
			tok.length = i - tok.offset
			tok.t = symbol(s)
			tokens = append(tokens, tok)
		case '"':
			tok := token{offset: i}
			s, err := lexString(&buf, s, &i)
			if err != nil {
				return nil, err
			}
			tok.t = s
			tok.length = i - tok.offset
			tokens = append(tokens, tok)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			tok := token{offset: i}
			n, err := lexNumber(s, &i)
			if err != nil {
				return nil, err
			}
			tok.t = n
			tok.length = i - tok.offset
			tokens = append(tokens, tok)
		default:
			if symstarter(s[i]) {
				tok := token{offset: i}
				buf.Reset()
				buf.WriteByte(s[i])
				i++
				for i < len(s) && symchar(s[i]) {
					buf.WriteByte(s[i])
					i++
				}
				s := buf.String()
				buf.Reset()
				switch s {
				case "null":
					tok.t = nil
				case "true":
					tok.t = true
				case "false":
					tok.t = false
				default:
					tok.t = symbol(s)
				}
				tok.length = i - tok.offset
				tokens = append(tokens, tok)
			} else {
				return nil, fmt.Errorf("unexpected '%v'", string(s[i]))
			}
		}
	}
}

func lexNumber(input string, pos *int) (int, error) {
	s := 1
	n := 0
	for *pos < len(input) {
		switch input[*pos] {
		case '-':
			s = -1
			*pos = *pos + 1
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			d := int(input[*pos]) - int('0')
			n = 10*n + d
			*pos = *pos + 1
		default:
			return s * n, nil
		}
	}
	return s * n, nil
}

func lexString(buf *bytes.Buffer, input string, pos *int) (string, error) {
	i := *pos
	i++
	buf.Reset()
	for {
		if i >= len(input) {
			return "", fmt.Errorf("unexpected end of input in string literal")
		}
		switch input[i] {
		case '\\':
			i++
			if i >= len(input) {
				return "", fmt.Errorf("unexpected end of input in string literal")
			}
			switch input[i] {
			case '"':
				i++
				buf.WriteByte('"')
			case 'b':
				i++
				buf.WriteByte('\b')
			case 'f':
				i++
				buf.WriteByte('\f')
			case 'n':
				i++
				buf.WriteByte('\n')
			case 'r':
				i++
				buf.WriteByte('\r')
			case 't':
				i++
				buf.WriteByte('\t')
			default:
				return "", fmt.Errorf("invalid string escape")
			}
		case '"':
			i++
			s := buf.String()
			buf.Reset()
			*pos = i
			return s, nil
		default:
			buf.WriteByte(input[i])
			i++
		}
	}
}

func symstarter(c byte) bool {
	switch {
	case c == '_':
		return true
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	default:
		return false
	}
}

func symchar(c byte) bool {
	switch c {
	case '_', '-', ':', '/', '.':
		return true
	default:
		switch {
		case c >= '0' && c <= '9':
			return true
		case c >= 'a' && c <= 'z':
			return true
		case c >= 'A' && c <= 'Z':
			return true
		default:
			return false
		}
	}
}
