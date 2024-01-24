package lang_lab

import (
	"fmt"
	"reflect"
)

type Program[C any] struct {
	root      astNode[C]
	variables map[string]astNode[C]
}

func (p Program[C]) Run(context C) {
	p.root.evaluate(context)
}

type Parser[C any] struct {
	lexer    Lexer
	language *Language[C]

	currToken token
	variables map[string]astNode[C]
}

func NewParser[C any](lexer Lexer, language *Language[C]) *Parser[C] {
	parser := &Parser[C]{
		lexer:     lexer,
		language:  language,
		variables: make(map[string]astNode[C]),
	}
	parser.advance()
	return parser
}

func (p *Parser[C]) advance() {
	p.currToken = p.lexer.nextToken()
}

// Parse runs the parser, returning either the root node of the AST or a parse error.
func (p *Parser[C]) Parse() (Program[C], error) {
	var statements []astNode[C]

parse:
	for {
		switch p.currToken.tpe {
		case tokenVariable:
			variableName := p.currToken.value
			node, err := p.parseExpression()
			if err != nil {
				return Program[C]{}, err
			}

			p.variables[variableName] = node

		case tokenLiteral:
			node, err := p.parseOperation()
			if err != nil {
				return Program[C]{}, err
			}
			statements = append(statements, node)

		case tokenLBracket, tokenRBracket, tokenInvalid:
			return Program[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered illegal token %s", p.currToken.value))

		case tokenEOF:
			break parse

		default:
		}

		p.advance()
	}

	return Program[C]{
		root: astNode[C]{
			returnType: nil,
			evaluate: func(context C) interface{} {
				for _, statement := range statements {
					statement.evaluate(context)
				}
				return nil
			},
		},
		variables: p.variables,
	}, nil
}

// parseExpression constructs an astNode to be assigned to a variable.
func (p *Parser[C]) parseExpression() (astNode[C], error) {
	p.advance()

	for {
		switch p.currToken.tpe {
		case tokenVariable:
			variable, exists := p.variables[p.currToken.value]
			if !exists {
				return astNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered undeclared variable %s", p.currToken.value))
			}
			return variable, nil

		case tokenLiteral:
			return p.parseOperation()

		case tokenLBracket:
			return p.parseList()

		case tokenEOF:
			return astNode[C]{}, fmtTokenErr(p.currToken, "unexpected end of expression")

		default:
			return astNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered illegal token %s", p.currToken.value))
		}
	}
}

// parseOperation constructs an astNode representing an operation in the given Language.
func (p *Parser[C]) parseOperation() (astNode[C], error) {
	operator := p.currToken

	p.advance()

	var operands []astNode[C]

	for {
		switch p.currToken.tpe {
		case tokenVariable:
			variable, exists := p.variables[p.currToken.value]
			if !exists {
				return astNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered undeclared variable %s", p.currToken.value))
			}
			operands = append(operands, variable)

		case tokenLiteral:
			node, err := p.language.parse(p.currToken, nil)
			if err != nil {
				return astNode[C]{}, err
			}
			operands = append(operands, node)

		case tokenLBracket:
			node, err := p.parseList()
			if err != nil {
				return astNode[C]{}, err
			}
			operands = append(operands, node)

		case tokenNewline, tokenEOF:
			node, err := p.language.parse(operator, operands)
			if err != nil {
				return astNode[C]{}, err
			}
			return node, nil

		default:
			return astNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered illegal token %s", p.currToken.value))
		}

		p.advance()
	}
}

// parseList constructs an astNode that constructs a list literal.
func (p *Parser[C]) parseList() (astNode[C], error) {
	var elementType reflect.Type
	var values []astNode[C]

	p.advance()

	for {
		switch p.currToken.tpe {
		case tokenLiteral:
			node, err := p.language.parse(p.currToken, nil)
			if err != nil {
				return astNode[C]{}, err
			}

			if elementType != nil && elementType != node.returnType {
				return astNode[C]{}, fmtTokenErr(p.currToken, "list must contain a single type")
			}

			if elementType == nil {
				elementType = node.returnType
			}

			values = append(values, node)

		case tokenRBracket:
			if elementType == nil {
				return astNode[C]{}, fmtTokenErr(p.currToken, "list must contain at least one element")
			}

			sliceType := reflect.SliceOf(elementType)
			return astNode[C]{
				returnType: sliceType,
				evaluate: func(context C) interface{} {
					result := reflect.MakeSlice(sliceType, 0, 0)
					for _, value := range values {
						result = reflect.Append(result, reflect.ValueOf(value.evaluate(context)))
					}
					return result.Interface()
				},
			}, nil

		case tokenEOF:
			return astNode[C]{}, fmtTokenErr(p.currToken, "unexpected end of list")

		default:
			return astNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("unexpected list element %s", p.currToken.value))
		}

		p.advance()
	}
}

// fmtTokenErr is used internally to return a message with line number.
func fmtTokenErr(t token, msg string) error {
	return fmt.Errorf("[line %d] %s", t.line, msg)
}
