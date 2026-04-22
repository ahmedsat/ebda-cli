package training

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenRegion TokenType = iota
	TokenFarm
	TokenCode
	TokenModule

	TokenNot
	TokenAnd
	TokenOr
	TokenXor

	TokenLParen
	TokenRParen
)

type Token struct {
	Type  TokenType
	Value string
}

func tokenize(input string) ([]Token, error) {
	var tokens []Token
	i := 0

	for i < len(input) {

		if unicode.IsSpace(rune(input[i])) {
			i++
			continue
		}

		switch input[i] {

		case '(':
			tokens = append(tokens, Token{Type: TokenLParen})
			i++

		case ')':
			tokens = append(tokens, Token{Type: TokenRParen})
			i++

		default:

			if unicode.IsLetter(rune(input[i])) {

				start := i
				for i < len(input) && (unicode.IsLetter(rune(input[i])) || input[i] == '_') {
					i++
				}

				word := strings.ToLower(input[start:i])

				switch word {

				case "and":
					tokens = append(tokens, Token{Type: TokenAnd})

				case "or":
					tokens = append(tokens, Token{Type: TokenOr})

				case "xor":
					tokens = append(tokens, Token{Type: TokenXor})

				case "not":
					tokens = append(tokens, Token{Type: TokenNot})

				case "region", "farm", "code", "module":

					for i < len(input) && unicode.IsSpace(rune(input[i])) {
						i++
					}

					if i >= len(input) || input[i] != '(' {
						return nil, fmt.Errorf("expected '(' after %s", word)
					}

					i++

					for i < len(input) && unicode.IsSpace(rune(input[i])) {
						i++
					}

					if i >= len(input) || input[i] != '"' {
						return nil, fmt.Errorf("expected string")
					}

					i++
					start = i

					for i < len(input) && input[i] != '"' {
						i++
					}

					if i >= len(input) {
						return nil, fmt.Errorf("unterminated string")
					}

					value := input[start:i]
					i++

					for i < len(input) && unicode.IsSpace(rune(input[i])) {
						i++
					}

					if i >= len(input) || input[i] != ')' {
						return nil, fmt.Errorf("expected ')'")
					}

					i++

					switch word {
					case "region":
						tokens = append(tokens, Token{Type: TokenRegion, Value: value})
					case "farm":
						tokens = append(tokens, Token{Type: TokenFarm, Value: value})
					case "code":
						tokens = append(tokens, Token{Type: TokenCode, Value: value})
					case "module":
						tokens = append(tokens, Token{Type: TokenModule, Value: value})
					}

				default:
					return nil, fmt.Errorf("unknown identifier %s", word)
				}

			} else {
				return nil, fmt.Errorf("unexpected character %c", input[i])
			}
		}
	}

	return tokens, nil
}

var precedence = map[TokenType]int{
	TokenOr:  1,
	TokenXor: 2,
	TokenAnd: 3,
	TokenNot: 4,
}

func Compile(input string) ([]Opt, error) {

	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}

	var output []Opt
	var stack []Token

	pushOp := func(t Token) {

		for len(stack) > 0 {

			top := stack[len(stack)-1]

			if top.Type == TokenLParen {
				break
			}

			if precedence[top.Type] < precedence[t.Type] {
				break
			}

			stack = stack[:len(stack)-1]
			output = append(output, tokenToOpt(top))
		}

		stack = append(stack, t)
	}

	for _, t := range tokens {

		switch t.Type {

		case TokenRegion, TokenFarm, TokenCode, TokenModule:
			output = append(output, tokenToOpt(t))

		case TokenNot, TokenAnd, TokenOr, TokenXor:
			pushOp(t)

		case TokenLParen:
			stack = append(stack, t)

		case TokenRParen:

			for {
				if len(stack) == 0 {
					return nil, fmt.Errorf("mismatched parentheses")
				}

				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if top.Type == TokenLParen {
					break
				}

				output = append(output, tokenToOpt(top))
			}
		}
	}

	for i := len(stack) - 1; i >= 0; i-- {

		if stack[i].Type == TokenLParen {
			return nil, fmt.Errorf("mismatched parentheses")
		}

		output = append(output, tokenToOpt(stack[i]))
	}

	return output, nil
}

func tokenToOpt(t Token) Opt {

	switch t.Type {

	case TokenRegion:
		return Opt{OptType: Region, Operand: t.Value}

	case TokenFarm:
		return Opt{OptType: Farm, Operand: t.Value}

	case TokenCode:
		return Opt{OptType: Code, Operand: t.Value}

	case TokenModule:
		return Opt{OptType: Module, Operand: t.Value}

	case TokenNot:
		return Opt{OptType: Not}

	case TokenAnd:
		return Opt{OptType: And}

	case TokenOr:
		return Opt{OptType: Or}

	case TokenXor:
		return Opt{OptType: Xor}
	}

	panic("invalid token")
}
