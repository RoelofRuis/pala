package pala

import (
	"fmt"
	"reflect"
)

type Parser[C any] struct {
	lexer            Lexer
	language         *Language[C]
	currToken        token
	program          Program[C]
	definedVariables map[string]reflect.Type
}

func NewParser[C any](lexer Lexer, language *Language[C]) *Parser[C] {
	parser := &Parser[C]{
		lexer:    lexer,
		language: language,
		program: Program[C]{
			variables: make(map[string]interface{}),
		},
		definedVariables: make(map[string]reflect.Type),
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
			variableToken := p.currToken

			expr, err := p.parseExpression()
			if err != nil {
				return Program[C]{}, err
			}

			node, err := p.writeVariable(variableToken, expr)
			if err != nil {
				return Program[C]{}, err
			}

			statements = append(statements, node)

		case tokenLiteral:
			node, err := p.parseOperation()
			if err != nil {
				return Program[C]{}, err
			}
			statements = append(statements, node)

		case tokenEOF:
			break parse

		case tokenComment, tokenNewline:

		default:
			return Program[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("encountered illegal token %s", p.currToken.value))
		}

		p.advance()
	}

	p.program.root = rootNode[C](statements)

	return p.program, nil
}

// parseExpression constructs an astNode to be assigned to a variable.
func (p *Parser[C]) parseExpression() (astNode[C], error) {
	p.advance()

	for {
		switch p.currToken.tpe {
		case tokenVariable:
			return p.readVariable(p.currToken)

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
	multiLine := false

	p.advance()

	var operands []astNode[C]

	for {
		switch p.currToken.tpe {
		case tokenLParen:
			if multiLine {
				return astNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("invalid additional opening parenthesis"))
			}
			multiLine = true

		case tokenRParen:
			if !multiLine {
				return astNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("invalid closing parenthesis"))
			}
			multiLine = false

		case tokenVariable:
			variable, err := p.readVariable(p.currToken)
			if err != nil {
				return astNode[C]{}, err
			}
			operands = append(operands, variable)

		case tokenLiteral:
			node, err := p.language.parseLiteral(p.currToken)
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

		case tokenNewline:
			if multiLine {
				break
			}
			node, err := p.language.parseOperation(operator, operands)
			if err != nil {
				return astNode[C]{}, err
			}
			return node, nil

		case tokenEOF:
			if multiLine {
				return astNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("missing closing parenthesis"))
			}
			node, err := p.language.parseOperation(operator, operands)
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
			node, err := p.language.parseLiteral(p.currToken)
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

		case tokenLBracket:
			node, err := p.parseList()
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
				return nilNode[C](), nil
			}

			sliceType := reflect.SliceOf(elementType)
			return sliceNode[C](sliceType, values), nil

		case tokenEOF, tokenNewline:
			return astNode[C]{}, fmtTokenErr(p.currToken, "unexpected end of list")

		default:
			return astNode[C]{}, fmtTokenErr(p.currToken, fmt.Sprintf("unexpected list element %s", p.currToken.value))
		}

		p.advance()
	}
}

// writeVariable writes a variable to the program variables.
func (p *Parser[C]) writeVariable(variableName token, value astNode[C]) (astNode[C], error) {
	p.definedVariables[variableName.value] = value.returnType

	return astNode[C]{
		returnType: nil,
		evaluate: func(context C) interface{} {
			p.program.variables[variableName.value] = value.evaluate(context)
			return nil
		},
	}, nil
}

// readVariable reads a variable from the program variables.
func (p *Parser[C]) readVariable(variableName token) (astNode[C], error) {
	varType, isDefined := p.definedVariables[variableName.value]
	if !isDefined {
		return astNode[C]{}, fmtTokenErr(variableName, fmt.Sprintf("encountered undeclared variable %s", variableName.value))
	}
	return astNode[C]{
		returnType: varType,
		evaluate: func(context C) interface{} {
			return p.program.variables[variableName.value]
		},
	}, nil
}

// fmtTokenErr is used internally to return a message with line number.
func fmtTokenErr(t token, msg string) error {
	return fmt.Errorf("[line %d] %s", t.line, msg)
}
