package lang_lab

import (
	"fmt"
	"reflect"
)

type Parser[C any] struct {
	lexer    Lexer
	language *Language[C]

	currToken token
	variables map[string]ASTNode[C]
}

func NewParser[C any](lexer Lexer, language *Language[C]) *Parser[C] {
	parser := &Parser[C]{
		lexer:     lexer,
		language:  language,
		variables: make(map[string]ASTNode[C]),
	}
	parser.advance()
	return parser
}

func (p *Parser[C]) advance() {
	p.currToken = p.lexer.nextToken()
}

// Parse runs the parser, returning either the root node of the AST or a parse error.
func (p *Parser[C]) Parse() (ASTNode[C], error) {
	var statements []ASTNode[C]

	for {
		switch p.currToken.tpe {
		case tokenVariable:
			variableName := p.currToken.value
			node, err := p.parseExpression()
			if err != nil {
				return ASTNode[C]{}, err
			}

			p.variables[variableName] = node

		case tokenLiteral:
			node, err := p.parseOperation()
			if err != nil {
				return ASTNode[C]{}, err
			}
			statements = append(statements, node)

		case tokenLBracket, tokenRBracket, tokenInvalid:
			return ASTNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered illegal token %s", p.currToken.value))

		case tokenEOF:
			return ASTNode[C]{
				returnType: nil,
				Evaluate: func(context C) interface{} {
					for _, statement := range statements {
						statement.Evaluate(context)
					}
					return nil
				},
			}, nil

		default:
		}

		p.advance()
	}
}

// parseExpression constructs an ASTNode to be assigned to a variable.
func (p *Parser[C]) parseExpression() (ASTNode[C], error) {
	p.advance()

	for {
		switch p.currToken.tpe {
		case tokenVariable:
			variable, exists := p.variables[p.currToken.value]
			if !exists {
				return ASTNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered undeclared variable %s", p.currToken.value))
			}
			return variable, nil

		case tokenLiteral:
			return p.parseOperation()

		case tokenLBracket:
			return p.parseList()

		case tokenEOF:
			return ASTNode[C]{}, fmtTokenErr(p.currToken, "unexpected end of expression")

		default:
			return ASTNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered illegal token %s", p.currToken.value))
		}
	}
}

// parseOperation constructs an ASTNode representing an operation in the given Language.
func (p *Parser[C]) parseOperation() (ASTNode[C], error) {
	operator := p.currToken

	p.advance()

	var operands []ASTNode[C]

	for {
		switch p.currToken.tpe {
		case tokenVariable:
			variable, exists := p.variables[p.currToken.value]
			if !exists {
				return ASTNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered undeclared variable %s", p.currToken.value))
			}
			operands = append(operands, variable)

		case tokenLiteral:
			node, err := p.language.parse(p.currToken, nil)
			if err != nil {
				return ASTNode[C]{}, err
			}
			operands = append(operands, node)

		case tokenLBracket:
			node, err := p.parseList()
			if err != nil {
				return ASTNode[C]{}, err
			}
			operands = append(operands, node)

		case tokenNewline, tokenEOF:
			node, err := p.language.parse(operator, operands)
			if err != nil {
				return ASTNode[C]{}, err
			}
			return node, nil

		default:
			return ASTNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered illegal token %s", p.currToken.value))
		}

		p.advance()
	}
}

// parseList constructs an ASTNode that constructs a list literal.
func (p *Parser[C]) parseList() (ASTNode[C], error) {
	var elementType reflect.Type
	var values []ASTNode[C]

	p.advance()

	for {
		switch p.currToken.tpe {
		case tokenLiteral:
			node, err := p.language.parse(p.currToken, nil)
			if err != nil {
				return ASTNode[C]{}, err
			}

			if elementType != nil && elementType != node.returnType {
				return ASTNode[C]{}, fmtTokenErr(p.currToken, "list must contain a single type")
			}

			if elementType == nil {
				elementType = node.returnType
			}

			values = append(values, node)

		case tokenRBracket:
			if elementType == nil {
				return ASTNode[C]{}, fmtTokenErr(p.currToken, "list must contain at least one element")
			}

			sliceType := reflect.SliceOf(elementType)
			return ASTNode[C]{
				returnType: sliceType,
				Evaluate: func(context C) interface{} {
					result := reflect.MakeSlice(sliceType, 0, 0)
					for _, value := range values {
						result = reflect.Append(result, reflect.ValueOf(value.Evaluate(context)))
					}
					return result.Interface()
				},
			}, nil

		case tokenEOF:
			return ASTNode[C]{}, fmtTokenErr(p.currToken, "unexpected end of list")

		default:
			return ASTNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("unexpected list element %s", p.currToken.value))
		}

		p.advance()
	}
}

// fmtTokenErr is used internally to return a message with line number.
func fmtTokenErr(t token, msg string) error {
	return fmt.Errorf("[line %d] %s", t.line, msg)
}
