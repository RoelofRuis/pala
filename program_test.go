package pala

import (
	"strings"
	"testing"
)

func Test_EvaluateImmediately(t *testing.T) {
	lang := NewLanguage[*context]()
	lang.BindOperator("+", plus)
	lang.BindLiteralEvaluator(parseInt)

	parser := NewParser(
		NewLexer(strings.NewReader("$a + 2 3")),
		lang,
	)

	prog, err := parser.Parse()
	if err != nil {
		t.Fatalf("expected program to be parsed:\n%s", err)
	}

	ctx := &context{}
	prog.Run(ctx)

	actualLog := ctx.String()
	expectedLog := "added 2 and 4"
	if actualLog != expectedLog {
		t.Errorf("expected log to contain '%s' but got '%s'", expectedLog, actualLog)
	}
}
