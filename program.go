package pala

import "reflect"

type Program[C any] struct {
	root      astNode[C]
	variables map[string]interface{}
}

func (p Program[C]) Run(context C) {
	p.root.evaluate(context)
}

type astNode[C any] struct {
	returnType reflect.Type
	evaluate   func(context C) interface{}
}

// rootNode creates an astNode that evaluates all statements and returns nil.
func rootNode[C any](statements []astNode[C]) astNode[C] {
	return astNode[C]{
		returnType: nil,
		evaluate: func(context C) interface{} {
			for _, statement := range statements {
				statement.evaluate(context)
			}
			return nil
		},
	}
}

// nilNode creates an astNode that evaluates to nil
func nilNode[C any]() astNode[C] {
	return valueNode[C](nil, nil)
}

// valueNode creates an astNode that evaluates to value.
func valueNode[C any](returnType reflect.Type, value interface{}) astNode[C] {
	return astNode[C]{
		returnType: returnType,
		evaluate:   func(context C) interface{} { return value },
	}
}

// sliceNode creates an astNode that evaluates to a slice of the given type.
func sliceNode[C any](returnType reflect.Type, values []astNode[C]) astNode[C] {
	return astNode[C]{
		returnType: returnType,
		evaluate: func(context C) interface{} {
			result := reflect.MakeSlice(returnType, 0, 0)
			for _, value := range values {
				result = reflect.Append(result, reflect.ValueOf(value.evaluate(context)))
			}
			return result.Interface()
		},
	}
}

// emptySliceNode creates an astNode that evaluates to an empty slice of the given type.
func emptySliceNode[C any](returnType reflect.Type) astNode[C] {
	return sliceNode[C](returnType, nil)
}

// operatorNode creates an astNode that evaluates the given operator with the given operands.
func operatorNode[C any](returnType reflect.Type, acceptsContext bool, operator reflect.Value, operands []astNode[C]) astNode[C] {
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

			result := operator.Call(arguments)
			if returnType == nil {
				return nil
			}
			return result[0].Interface()
		},
	}
}
