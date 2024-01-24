package pala

type Program[C any] struct {
	root      astNode[C]
	variables map[string]astNode[C]
}

func (p Program[C]) Run(context C) {
	p.root.evaluate(context)
}
