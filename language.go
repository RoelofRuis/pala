package pala

import (
	"fmt"
	"reflect"
)

// Language contains evaluators that convert string symbols to the appropriate literals and functions.
// You construct the language by defining available literals and operations using the BindLiteralEvaluator and
// BindOperator methods.
type Language[C any] struct {
	operators map[string]func(operands []astNode[C]) (astNode[C], error)
	literals  []func(token token) (astNode[C], error)
}

// NewLanguage constructs an empty Language.
func NewLanguage[C any]() *Language[C] {
	return &Language[C]{
		operators: make(map[string]func(operands []astNode[C]) (astNode[C], error)),
		literals:  []func(token token) (astNode[C], error){},
	}
}

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
		return valueNode[C](returnType, value), nil
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
	numExpectedOperands := funcType.NumIn()
	for i := 0; i < funcType.NumIn(); i++ {
		if i == 0 && funcType.In(0) == reflect.TypeOf(zero).Elem() {
			acceptsContext = true
			numExpectedOperands -= 1
		} else {
			argTypes = append(argTypes, funcType.In(i))
		}
	}

	var returnType reflect.Type
	if funcType.NumOut() == 1 {
		returnType = funcType.Out(0)
	}

	operator := func(operands []astNode[C]) (astNode[C], error) {
		if numExpectedOperands != len(operands) {
			return astNode[C]{}, fmt.Errorf("operator %s expected %d operands but got %d", symbol, numExpectedOperands, len(operands))
		}
		for i, operand := range operands {
			if argTypes[i].Kind() == reflect.Slice && operand.returnType == nil {
				// slice types accept nil: this equates to an empty slice of the appropriate type.
				operands[i] = emptySliceNode[C](argTypes[i])
				continue
			}
			if argTypes[i] == operand.returnType {
				continue
			}
			if argTypes[i].Kind() == reflect.Interface && operand.returnType.Implements(argTypes[i]) {
				continue
			}

			return astNode[C]{}, fmt.Errorf("operand %d of operator %s expects %s but got %v", i, symbol, argTypes[i], operand.returnType)
		}
		return operatorNode[C](returnType, acceptsContext, funcValue, operands), nil
	}

	l.operators[symbol] = operator
}

func (l *Language[C]) parseLiteral(token token) (astNode[C], error) {
	for _, literal := range l.literals {
		node, err := literal(token)
		if err != nil {
			continue
		}
		return node, nil
	}
	return astNode[C]{}, fmt.Errorf("unknown literal %s", token.value)
}

func (l *Language[C]) parseOperation(token token, operands []astNode[C]) (astNode[C], error) {
	operator, has := l.operators[token.value]
	if !has {
		return astNode[C]{}, fmt.Errorf("unknown operator %s", token.value)
	}
	return operator(operands)
}

var stringType = reflect.TypeOf("")
