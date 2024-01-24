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
