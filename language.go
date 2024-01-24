package pala

import (
	"fmt"
	"reflect"
)

type Language[C any] struct {
	operators map[string]func(operands []astNode[C]) (astNode[C], error)
	literals  []func(token token) (astNode[C], error)
}

func NewLanguage[C any]() *Language[C] {
	return &Language[C]{
		operators: make(map[string]func(operands []astNode[C]) (astNode[C], error)),
		literals:  []func(token token) (astNode[C], error){},
	}
}

var stringType = reflect.TypeOf("")

// BindLiteralEvaluator adds a literal evaluator to the language.
// It must be provided with a function with signature `func(string) (any, error)`. This function should try to parse the
// given string into a literal and return it. It may fail with an error, in which case the parser will proceed to the
// next literal evaluator function that was bound.
func (l *Language[C]) BindLiteralEvaluator(evaluator interface{}) {
	funcValue := reflect.ValueOf(evaluator)

	if funcValue.Kind() != reflect.Func {
		panic("function is required")
	}

	funcType := funcValue.Type()

	if funcType.NumOut() != 2 ||
		funcType.NumIn() != 1 ||
		funcType.In(0) != stringType {
		panic("function must have signature func(string) (any, error)")
	}

	returnType := funcType.Out(0)

	primitive := func(token token) (astNode[C], error) {
		values := funcValue.Call([]reflect.Value{reflect.ValueOf(token.value)})
		value := values[0].Interface()
		err := values[1].Interface()
		if err != nil {
			return astNode[C]{}, err.(error)
		}
		return astNode[C]{
			returnType: returnType,
			evaluate: func(context C) interface{} {
				return value
			},
		}, nil
	}

	l.literals = append(l.literals, primitive)
}

// BindOperator binds an operator constructing function to be triggered when the given symbol is encountered.
// The constructor function can have any number of input arguments of any type, and can only have one or zero return
// values.
// If the first value is of the context type `C` of the language, the context will be passed to it during
// interpretation.
func (l *Language[C]) BindOperator(symbol string, constructor interface{}) {
	funcValue := reflect.ValueOf(constructor)

	if funcValue.Kind() != reflect.Func {
		panic("function is required")
	}

	funcType := funcValue.Type()

	if funcType.NumOut() > 1 {
		panic("functions must have zero or one return values")
	}

	var zero [0]C
	var argTypes []reflect.Type
	acceptsContext := false
	for i := 0; i < funcType.NumIn(); i++ {
		if i == 0 && funcType.In(0) == reflect.TypeOf(zero).Elem() {
			acceptsContext = true
		} else {
			argTypes = append(argTypes, funcType.In(i))
		}
	}

	var returnType reflect.Type
	if funcType.NumOut() == 1 {
		returnType = funcType.Out(0)
	}

	operator := func(operands []astNode[C]) (astNode[C], error) {
		for i, operand := range operands {
			if argTypes[i] != operand.returnType {
				return astNode[C]{}, fmt.Errorf("operator operand %d expects %s but got %s", i, argTypes[i], operand.returnType)
			}
		}
		return astNode[C]{
			returnType: returnType,
			evaluate: func(context C) interface{} {
				var arguments []reflect.Value
				if acceptsContext {
					arguments = append(arguments, reflect.ValueOf(context))
				}

				for _, operand := range operands {
					arguments = append(arguments, reflect.ValueOf(operand.evaluate(context)))
				}

				result := funcValue.Call(arguments)
				if returnType == nil {
					return nil
				}
				return result[0].Interface()
			},
		}, nil
	}

	l.operators[symbol] = operator
}

func (l *Language[C]) parse(token token, operands []astNode[C]) (astNode[C], error) {
	if operands == nil {
		for _, literal := range l.literals {
			node, err := literal(token)
			if err != nil {
				continue
			}
			return node, nil
		}
		return astNode[C]{}, fmt.Errorf("unknown literal %s", token.value)
	}

	for symbol, builder := range l.operators {
		if symbol == token.value {
			return builder(operands)
		}
	}

	return astNode[C]{}, fmt.Errorf("unknown operator %s", token.value)
}
