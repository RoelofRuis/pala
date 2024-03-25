package pala

import (
	"strings"
	"testing"
)

func Test_ParseAndRun(t *testing.T) {
	tests := []struct {
		name        string
		program     string
		expectedLog string
	}{
		{
			"empty program",
			"",
			"",
		},
		{
			"comment",
			"# this does nothing",
			"",
		},
		{
			"nullary operator",
			"dbg",
			"debug",
		},
		{
			"unary operator",
			"neg 4",
			"negated 4",
		},
		{
			"binary operator",
			"+ 2 3",
			"added 2 and 3",
		},
		{
			"binary operator with additional spaces",
			"+  2      3",
			"added 2 and 3",
		},
		{
			"binary operator save to variable",
			"$a + 2 3",
			"added 2 and 3",
		},
		{
			"multiple statements",
			"$a + 2 3\n$b * $a 5",
			"added 2 and 3\nmultiplied 5 and 5",
		},
		{
			"working with lists",
			"min [1 2 3]",
			"finding min of [1,2,3]",
		},
		{
			"allow any implementer as argument",
			"echo \"string\"\necho 4",
			"",
		},
		{
			"nested lists",
			"$a shortest [[1 2 3 4][5 6 7][8 9]]\nmin $a",
			"finding min of [8,9]",
		},
		{
			"empty list",
			"$a min []\n",
			"finding min of []",
		},
		{
			"statement split on multiple lines",
			"+ (\n    4\n    5\n)",
			"added 4 and 5",
		},
		{
			"empty lines between statements",
			"+ 3 4\n\nmin [2 3]",
			"added 3 and 4\nfinding min of [2,3]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang := NewLanguage[*context]()
			lang.BindOperator("dbg", debug)
			lang.BindOperator("neg", neg)
			lang.BindOperator("min", smallest)
			lang.BindOperator("echo", echo)
			lang.BindOperator("+", plus)
			lang.BindOperator("*", mul)
			lang.BindOperator("shortest", shortest)
			lang.BindLiteralEvaluator(ParseInt)
			lang.BindLiteralEvaluator(ParseQuotedString)

			parser := NewParser(
				NewLexer(strings.NewReader(tt.program)),
				lang,
			)

			prog, err := parser.Parse()
			if err != nil {
				t.Fatalf("expected program to be parsed:\n%s", err)
			}

			ctx := &context{}
			prog.Run(ctx)

			actualLog := ctx.String()
			if actualLog != tt.expectedLog {
				t.Errorf("expected log to contain '%s' but got '%s'", tt.expectedLog, actualLog)
			}
		})
	}
}
